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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-grok"

	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/filter"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/context_values"
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

	client *s3.Client
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
		slog.Error("Error getting S3 client", "error", err)
		return err
	}
	s.client = client

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
		// fatal error - log and return
		slog.Error("error adding grok patterns", "error", err)
		return fmt.Errorf("error adding grok patterns: %v", err)
	}

	if s.Config.Prefix != nil {
		prefix = *s.Config.Prefix
		var newOptionalLayouts []string
		for _, l := range optionalLayouts {
			newOptionalLayouts = append(newOptionalLayouts, fmt.Sprintf("%s%s", prefix, l))
		}
		// Add support for collecting logs from S3 buckets that use a flat structure (i.e., without directory-style prefixes).
		// Currently, if a prefix is specified in the config, it is prepended to the layout pattern.
		// For example, if the prefix is "2025-06-06" and the layout is "%{YEAR:year}-%{MONTHNUM:month}-%{MONTHDAY:day}-%{HOUR:hour}-%{MINUTE:minute}-%{SECOND:second}-%{DATA:suffix}",
		// the resulting layout becomes "2025-06-06%{YEAR:year}-%{MONTHNUM:month}-%{MONTHDAY:day}-%{HOUR:hour}-%{MINUTE:minute}-%{SECOND:second}-%{DATA:suffix}",
		// which breaks log collection from buckets using a flat file structure.
		// To address this, we're preserving the existing behavior for directory-style buckets,
		// while adding support for flat buckets as a new, optional configuration path.
		optionalLayouts = append(optionalLayouts, newOptionalLayouts...)
	}

	// walkS3 should only return fatal errors
	err = s.walkS3(ctx, s.Config.Bucket, prefix, optionalLayouts, filterMap, g)
	if err != nil {
		slog.Error("error walking S3 bucket", "bucket", s.Config.Bucket, "error", err)
		return fmt.Errorf("%s: %s", s.Config.Bucket, err.Error())
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
		slog.Error("failed to download artifact", "bucket", s.Config.Bucket, "key", info.Name, "error", err)
		return fmt.Errorf("%s: failed to download artifact from %s", info.Name, s.Config.Bucket)
	}
	defer getObjectOutput.Body.Close()

	// Get the size of the object
	size := typehelpers.Int64Value(getObjectOutput.ContentLength)

	// copy the object data to a temp file
	localFilePath := path.Join(s.TempDir, info.Name)
	localFileDir := filepath.Dir(localFilePath)

	// ensure the directory exists of the file to write to
	if err := os.MkdirAll(localFileDir, 0755); err != nil {
		slog.Error("failed to create directory", "bucket", s.Config.Bucket, "key", info.Name, "dir", localFileDir, "error", err)
		return fmt.Errorf("%s: failed to download artifact from %s", info.Name, s.Config.Bucket)
	}

	// Create a local file to write the data to
	outFile, err := os.Create(localFilePath)
	if err != nil {
		slog.Error("failed to create file", "bucket", s.Config.Bucket, "key", info.Name, "file", outFile, "error", err)
		return fmt.Errorf("%s: failed to download artifact from %s", info.Name, s.Config.Bucket)
	}
	defer outFile.Close()

	// Write the data to the local file
	_, err = io.Copy(outFile, getObjectOutput.Body)
	if err != nil {
		slog.Error("failed to write file content", "bucket", s.Config.Bucket, "key", info.Name, "file", outFile, "error", err)
		return fmt.Errorf("%s: failed to download artifact from %s", info.Name, s.Config.Bucket)
	}

	// notify observers of the downloaded artifact
	return s.OnArtifactDownloaded(ctx, types.NewDownloadedArtifactInfo(info, localFilePath, size))
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
	executionId, err := context_values.ExecutionIdFromContext(ctx)
	if err != nil {
		return err
	}

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			// fatal error - log and return
			slog.Error("error getting next page", "bucket", bucket, "prefix", prefix, "error", err)
			return fmt.Errorf("error getting next page, %w", err)
		}

		// Directories
		for _, dir := range page.CommonPrefixes {
			dirPrefix := typehelpers.SafeString(*dir.Prefix)
			err = s.WalkNode(ctx, dirPrefix, "", layouts, true, g, filterMap)
			if err != nil {
				// ignore skip dir error as this means directory isn't one we want to dive into
				if errors.Is(err, fs.SkipDir) {
					continue
				}
				// non-fatal error - log and notify
				slog.Error("error obtaining directory info", "key", dirPrefix, "error", err)
				s.NotifyError(ctx, executionId, fmt.Errorf("%s: failed to obtain directory info", dirPrefix))
				err = nil
				continue
			}
			err = s.walkS3(ctx, bucket, dirPrefix, layouts, filterMap, g)
			if err != nil {
				// non-fatal error - log and notify
				slog.Error("error walking S3 bucket", "bucket", bucket, "prefix", dirPrefix, "error", err)
				s.NotifyError(ctx, executionId, fmt.Errorf("%s: %s", bucket, err.Error()))
				err = nil
			}
		}

		// Files
		for _, obj := range page.Contents {
			objKey := typehelpers.SafeString(*obj.Key)
			if objKey == "" {
				slog.Debug("skipping empty object key")
				continue
			}
			err = s.WalkNode(ctx, objKey, "", layouts, false, g, filterMap)
			if err != nil {
				// non-fatal error - log and notify
				slog.Error("error obtaining artifact info", "key", objKey, "error", err)
				s.NotifyError(ctx, executionId, fmt.Errorf("%s: failed to obtain artifact info", objKey))
				err = nil
			}
		}
	}

	return nil
}
