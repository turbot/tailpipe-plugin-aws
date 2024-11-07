package tables

// Package tables implements the AWS ALB (Application Load Balancer) access log table.
// This implementation handles parsing and structuring ALB access logs into queryable data.

import (
	"context"
	"github.com/turbot/tailpipe-plugin-aws/mappers"

	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type AlbAccessLogTable struct {
	table.TableImpl[*rows.AlbAccessLog, *AlbAccessLogTableConfig, *config.AwsConnection]
}

func NewAlbAccessLogTable() table.Table {
	return &AlbAccessLogTable{}
}

func (t *AlbAccessLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	if err := t.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	// Set the mapper
	t.Mapper = mappers.NewAlbAccessLogMapper()
	return nil
}

func (t *AlbAccessLogTable) Identifier() string {
	return "aws_alb_access_log"
}

func (t *AlbAccessLogTable) GetRowSchema() any {
	return &rows.AlbAccessLog{}
}

func (t *AlbAccessLogTable) GetConfigSchema() parse.Config {
	return &AlbAccessLogTableConfig{}
}

func (t *AlbAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (t *AlbAccessLogTable) EnrichRow(row *rows.AlbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AlbAccessLog, error) {
	if err := row.EnrichRow(sourceEnrichmentFields); err != nil {
		return nil, err
	}
	return row, nil
}
