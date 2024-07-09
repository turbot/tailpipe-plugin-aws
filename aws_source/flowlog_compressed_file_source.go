package aws_source

import (
	"bufio"
	"compress/gzip"
	"context"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/source"
	"log/slog"
	"os"
	"path/filepath"
)

// FlowlogCompressedFileSource source is responsible for collecting audit logs from Turbot Pipes API
type FlowlogCompressedFileSource struct {
	source.Base
	Config CompressedFileSourceConfig
}

func (c *FlowlogCompressedFileSource) Identifier() string {
	return "aws_compressed_file_source"
}

func NewFlowlogCompressedFileSource(config CompressedFileSourceConfig) plugin.Source {
	return &FlowlogCompressedFileSource{
		Config: config,
	}
}

// TODO combine with CloudtrailCompressedFileSource
func (c *FlowlogCompressedFileSource) Collect(ctx context.Context, req *proto.CollectRequest) error {
	// tactical
	//List all gz files in each path directory and call ExtractArtifactRows for each
	for _, path := range c.Config.Paths {
		// list gz files on path
		// TODO recursive or not?
		files, err := findGzFiles(path)
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

// TODO tactical - sort out file source
func findGzFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gz" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (c *FlowlogCompressedFileSource) ExtractArtifactRows(ctx context.Context, req *proto.CollectRequest, inputPath string) error {
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

	scanner := bufio.NewScanner(gzReader)
	for scanner.Scan() {
		// check context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// populate enrichment fields the the source is aware of
		// - in this case the source location
		sourceEnrichmentFields := &enrichment.CommonFields{
			TpSourceLocation: &inputPath,
		}

		line := scanner.Text()
		// call base OnRow
		c.OnRow(req, line, sourceEnrichmentFields)
	}
	return scanner.Err()
}
