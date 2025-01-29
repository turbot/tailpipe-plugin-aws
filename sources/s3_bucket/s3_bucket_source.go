package s3_bucket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-grok"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/filter"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const (
	AwsS3BucketSourceIdentifier = "aws_s3_bucket"
	defaultBucketRegion         = "us-east-1"
)

// AwsS3BucketSource is a [ArtifactSource] implementation that reads artifacts from an S3 bucket
type AwsS3BucketSource struct {
	artifact_source.ArtifactSourceImpl[*AwsS3BucketSourceConfig, *config.AwsConnection]

	client    *s3.Client
	errorList []error
}

func (s *AwsS3BucketSource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	slog.Info("Initializing AwsS3BucketSource")

	// call base to parse config and apply options
	if err := s.ArtifactSourceImpl.Init(ctx, params, opts...); err != nil {
		return err
	}

	// initialize client
	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}
	s.client = client

	s.errorList = []error{}

	slog.Info("Initialized AwsS3BucketSource", "bucket", s.Config.Bucket, "layout", s.Config.FileLayout)

	return nil
}

func (s *AwsS3BucketSource) Identifier() string {
	return AwsS3BucketSourceIdentifier
}

func (s *AwsS3BucketSource) Close() error {
	_ = os.RemoveAll(s.TempDir)
	return nil
}

func (s *AwsS3BucketSource) ValidateConfig() error {
	if s.Config.Bucket == "" {
		return fmt.Errorf("bucket is required and cannot be empty")
	}

	return nil
}

func (s *AwsS3BucketSource) DiscoverArtifacts(ctx context.Context) error {
	var prefix string
	layout := typehelpers.SafeString(s.Config.GetFileLayout())
	// if there are any optional segments, we expand them into all possible alternatives
	optionalLayouts := artifact_source.ExpandPatternIntoOptionalAlternatives(layout)

	filterMap := make(map[string]*filter.SqlFilter)

	g := grok.New()
	// add any patterns defined in config
	err := g.AddPatterns(s.Config.GetPatterns())
	if err != nil {
		return fmt.Errorf("error adding grok patterns: %v", err)
	}

	if s.Config.Prefix != nil {
		prefix = *s.Config.Prefix
		if !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		var newOptionalLayouts []string
		for _, l := range optionalLayouts {
			newOptionalLayouts = append(newOptionalLayouts, fmt.Sprintf("%s%s", prefix, l))
		}
		optionalLayouts = newOptionalLayouts
	}

	err = s.walkS3(ctx, s.Config.Bucket, prefix, optionalLayouts, filterMap, g)
	if err != nil {
		s.errorList = append(s.errorList, fmt.Errorf("error discovering artifacts in S3 bucket %s, %w", s.Config.Bucket, err))
	}

	if len(s.errorList) > 0 {
		return errors.Join(s.errorList...)
	}
	return nil
}

func (s *AwsS3BucketSource) DownloadArtifact(ctx context.Context, info *types.ArtifactInfo) error {
	// Get the object from S3
	getObjectOutput, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Config.Bucket,
		Key:    &info.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to download artifact, %w", err)
	}
	defer getObjectOutput.Body.Close()

	// Get the size of the object
	size := typehelpers.Int64Value(getObjectOutput.ContentLength)

	// copy the object data to a temp file
	localFilePath := path.Join(s.TempDir, info.Name)
	// ensure the directory exists of the file to write to
	if err := os.MkdirAll(filepath.Dir(localFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for file, %w", err)
	}

	// Create a local file to write the data to
	outFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file, %w", err)
	}
	defer outFile.Close()

	// Write the data to the local file
	_, err = io.Copy(outFile, getObjectOutput.Body)
	if err != nil {
		return fmt.Errorf("failed to write data to file, %w", err)
	}

	// notify observers of the downloaded artifact
	downloadInfo := types.NewDownloadedArtifactInfo(info, localFilePath, size)

	return s.OnArtifactDownloaded(ctx, downloadInfo)
}

func (s *AwsS3BucketSource) getClient(ctx context.Context) (*s3.Client, error) {
	// get the client configuration
	tempRegion := defaultBucketRegion
	cfg, err := s.Connection.GetClientConfiguration(ctx, &tempRegion)
	if err != nil {
		return nil, fmt.Errorf("unable to get client configuration, %w", err)
	}

	region, err := manager.GetBucketRegion(ctx, s3.NewFromConfig(*cfg), s.Config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("unable to get bucket region, %w", err)
	}

	cfg.Region = region

	if s.Connection.S3ForcePathStyle != nil {
		return s3.NewFromConfig(*cfg, func(o *s3.Options) {
			o.UsePathStyle = *s.Connection.S3ForcePathStyle
		}), nil
	}

	return s3.NewFromConfig(*cfg), nil
}

func (s *AwsS3BucketSource) walkS3(ctx context.Context, bucket string, prefix string, layouts []string, filterMap map[string]*filter.SqlFilter, g *grok.Grok) error {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page, %w", err)
		}

		// Directories
		for _, dir := range page.CommonPrefixes {
			err = s.WalkNode(ctx, *dir.Prefix, "", layouts, true, g, filterMap)
			if err != nil {
				if errors.Is(err, fs.SkipDir) {
					continue
				}
				return fmt.Errorf("error walking node, %w", err)
			}
			err = s.walkS3(ctx, bucket, *dir.Prefix, layouts, filterMap, g)
			if err != nil {
				s.errorList = append(s.errorList, err)
			}
		}

		// Files
		for _, obj := range page.Contents {
			err = s.WalkNode(ctx, *obj.Key, "", layouts, false, g, filterMap)
			if err != nil {
				s.errorList = append(s.errorList, fmt.Errorf("error parsing object %s, %w", *obj.Key, err))
			}
		}
	}

	return nil
}
