package aws_source

import (
	"context"
	"errors"
	"fmt"
	"strings"
	//sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	//"github.com/turbot/tailpipe-plugin-sdk/source"
)

type AwsS3BucketSourceConfig struct {
	Bucket     string   `json:"bucket"`
	Extensions []string `json:"extensions"`
}

type AwsS3BucketSource struct {
	Config AwsS3BucketSourceConfig

	//ctx            context.Context
	//observers      []source.SourceObserver
	//observersMutex sync.RWMutex
}

func (s *AwsS3BucketSource) Identifier() string {
	return "aws_s3_bucket"
}

//func (s *AwsS3BucketSource) Init(ctx context.Context) error {
//	s.ctx = ctx
//	return s.Validate()
//}
//
//func (s *AwsS3BucketSource) Context() context.Context {
//	return s.ctx
//}
//
//func (s *AwsS3BucketSource) Validate() error {
//	return nil
//}

//func (s *AwsS3BucketSource) AddObserver(observer source.SourceObserver) {
//	s.observersMutex.Lock()
//	defer s.observersMutex.Unlock()
//	s.observers = append(s.observers, observer)
//}
//
//func (s *AwsS3BucketSource) RemoveObserver(observer source.SourceObserver) {
//	s.observersMutex.Lock()
//	defer s.observersMutex.Unlock()
//	for i, o := range s.observers {
//		if o == observer {
//			s.observers = append(s.observers[:i], s.observers[i+1:]...)
//			break
//		}
//	}
//}
//
//func (s *AwsS3BucketSource) LoadConfig(configRaw []byte) error {
//	return sdkconfig.Load(configRaw, &s.Config)
//}

func (s *AwsS3BucketSource) ValidateConfig() error {
	if s.Config.Bucket == "" {
		return errors.New("bucket is required")
	}

	// Check the bucket exists
	// TODO
	/*
		_, err := os.Stat(s.Config.Path)
		if err != nil {
			return err
		}
	*/

	// Check format of extensions
	invalidExtensions := []string{}
	for _, e := range s.Config.Extensions {
		if len(e) == 0 {
			invalidExtensions = append(invalidExtensions, "<empty>")
		} else if e[0] != '.' {
			invalidExtensions = append(invalidExtensions, e)
		}
	}
	if len(invalidExtensions) > 0 {
		return fmt.Errorf("invalid extensions: %s", strings.Join(invalidExtensions, ","))
	}

	return nil
}

func (s *AwsS3BucketSource) DiscoverArtifacts(ctx context.Context) error {
	// TODO implement
	//
	//cfg, err := config.LoadDefaultConfig(ctx)
	//if err != nil {
	//	return fmt.Errorf("unable to load SDK config, %w", err)
	//}
	//
	//s3Client := s3.NewFromConfig(cfg)
	//paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
	//	Bucket: &s.Config.Bucket,
	//})
	//
	//for paginator.HasMorePages() {
	//	output, err := paginator.NextPage(ctx)
	//	if err != nil {
	//		return fmt.Errorf("failed to get page, %w", err)
	//	}
	//
	//	for _, object := range output.Contents {
	//		if util.IsValidExtension(*object.Key, s.Config.Extensions) {
	//			for _, observer := range s.observers {
	//				observer.NotifyArtifactDiscovered(&source.ArtifactInfo{Name: *object.Key})
	//			}
	//		}
	//	}
	//}

	return nil
}

func (s *AwsS3BucketSource) DownloadArtifact(ctx context.Context, ai any /**source.ArtifactInfo*/) error {
	//TODO implement
	//cfg, err := config.LoadDefaultConfig(ctx)
	//if err != nil {
	//	return fmt.Errorf("unable to load SDK config, %w", err)
	//}
	//
	//s3Client := s3.NewFromConfig(cfg)
	//getObjectOutput, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
	//	Bucket: &s.Config.Bucket,
	//	Key:    &ai.Name,
	//})
	//if err != nil {
	//	return fmt.Errorf("failed to download artifact, %w", err)
	//}
	//defer getObjectOutput.Body.Close()
	//
	//data, err := io.ReadAll(getObjectOutput.Body)
	//if err != nil {
	//	return fmt.Errorf("failed to read artifact data, %w", err)
	//}
	//
	//for _, observer := range s.observers {
	//	observer.NotifyArtifactDownloaded(&source.Artifact{ArtifactInfo: *ai, Data: data})
	//}

	return nil
}
