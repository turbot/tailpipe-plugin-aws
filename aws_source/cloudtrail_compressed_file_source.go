package aws_source

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/source"
	"log/slog"
	"os"
	"path/filepath"
)

// CloudtrailCompressedFileSource source is responsible for collecting audit logs from Turbot Pipes API
type CloudtrailCompressedFileSource struct {
	source.Base
	Config CompressedFileSourceConfig
}

func (c *CloudtrailCompressedFileSource) Identifier() string {
	return "aws_compressed_file_source"
}

func NewCloudtrailCompressedFileSource(config CompressedFileSourceConfig) plugin.Source {
	return &CloudtrailCompressedFileSource{
		Config: config,
	}
}

func (c *CloudtrailCompressedFileSource) Collect(ctx context.Context, req *proto.CollectRequest) error {
	// tactical
	//List all gz files in each path directory and call ExtractArtifactRows for each
	for _, path := range c.Config.Paths {
		// list gz files on path
		files, err := filepath.Glob(filepath.Join(path, "*.gz"))
		if err != nil {
			return err
		}

		for _, file := range files {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				slog.Debug("Processing file", "file", file)
				// Call ExtractArtifactRows for each gz file
				if err := c.ExtractArtifactRows(ctx, req, file); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *CloudtrailCompressedFileSource) ExtractArtifactRows(ctx context.Context, req *proto.CollectRequest, inputPath string) error {
	gzFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	var log aws_types.AWSCloudTrailBatch
	if err := json.NewDecoder(gzReader).Decode(&log); err != nil {
		return err
	}

	for _, record := range log.Records {
		// check context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// populate enrichment fields the the source is aware of
		// - in this case the source location
		sourceEnrichmentFields := &enrichment.CommonFields{
			TpSourceLocation: &inputPath,
		}

		// call base OnRow
		c.OnRow(req, record, sourceEnrichmentFields)
	}

	return nil

}
