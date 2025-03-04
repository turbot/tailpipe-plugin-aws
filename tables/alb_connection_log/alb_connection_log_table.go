package alb_connection_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AlbConnectionLogTableIdentifier = "aws_alb_connection_log"

const connectionLogFormat = `$timestamp $client_ip $client_port $listener_port $tls_protocol $tls_cipher $tls_handshake_latency "$leaf_client_cert_subject" $leaf_client_cert_validity $leaf_client_cert_serial_number $tls_verify_status $conn_trace_id`

type AlbConnectionLogTable struct{}

func (c *AlbConnectionLogTable) Identifier() string {
	return AlbConnectionLogTableIdentifier
}

func (c *AlbConnectionLogTable) GetSourceMetadata() []*table.SourceMetadata[*AlbConnectionLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.[^_]+_[^_]+_[^_]+_[^.]+.log.gz"),
	}

	return []*table.SourceMetadata[*AlbConnectionLog]{
		{
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*AlbConnectionLog](connectionLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*AlbConnectionLog](connectionLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *AlbConnectionLogTable) EnrichRow(row *AlbConnectionLog, sourceEnrichmentFields schema.SourceEnrichment) (*AlbConnectionLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	row.TpSourceIP = &row.ClientIP
	row.TpIps = append(row.TpIps, row.ClientIP)
	row.TpIndex = *row.ConnTraceID
	return row, nil
}

// GetDescription returns a description of the connection log table.
func (c *AlbConnectionLogTable) GetDescription() string {
	return "AWS ALB Connection logs capture detailed information about connection attempts to an Application Load Balancer, including TLS handshake details, client certificate data, and connection traceability identifiers."
}