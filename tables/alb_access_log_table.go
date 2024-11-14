package tables

// Package tables implements the AWS ALB (Application Load Balancer) access log table.
// This implementation handles parsing and structuring ALB access logs into queryable data.

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.AlbAccessLog, *AlbAccessLogTable]()
}

type AlbAccessLogTable struct {
	table.TableImpl[*rows.AlbAccessLog, *AlbAccessLogTableConfig, *config.AwsConnection]
}

func (t *AlbAccessLogTable) SupportedSource() []*table.SourceMetadata[*rows.AlbAccessLog] {
	return []*table.SourceMetadata[*rows.AlbAccessLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: mappers.NewAlbAccessLogMapper,
			Options:    []row_source.RowSourceOption{artifact_source.WithRowPerLine()},
		},
	}
}

func (t *AlbAccessLogTable) Identifier() string {
	return "aws_alb_access_log"
}

func (*AlbAccessLogTable) EnrichRow(row *rows.AlbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AlbAccessLog, error) {
	// Add source enrichment fields if provided
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Standard record enrichment
	row.TpID = xid.New().String()
	row.TpTimestamp = row.Timestamp
	row.TpIngestTimestamp = time.Now()
	// truncate timestamp to date
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	// Use ALB name as the index
	row.TpIndex = row.AlbName

	// IP-related enrichment
	if row.ClientIP != "" {
		row.TpSourceIP = &row.ClientIP
		row.TpIps = append(row.TpIps, row.ClientIP)
	}
	if row.TargetIP != nil {
		row.TpDestinationIP = row.TargetIP
		row.TpIps = append(row.TpIps, *row.TargetIP)
	}
	// Domain enrichment
	if row.DomainName != "" {
		row.TpDomains = append(row.TpDomains, row.DomainName)
	}

	// AWS resource linking
	if row.TargetGroupArn != "" {
		row.TpAkas = append(row.TpAkas, row.TargetGroupArn)
	}
	return row, nil
}
