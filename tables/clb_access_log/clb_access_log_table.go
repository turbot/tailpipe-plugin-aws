package clb_access_log

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

const ClbAccessLogTableIdentifier = "aws_clb_access_log"

const clbLogFormat = `$timestamp $elb $client_ip:$client_port $backend $request_processing_time $backend_processing_time $response_processing_time $elb_status_code $backend_status_code $received_bytes $sent_bytes "$request_http_method $request_url $request_http_version" "$user_agent" $ssl_cipher $ssl_protocol`

type ClbAccessLogTable struct{}

func (c *ClbAccessLogTable) Identifier() string {
	return ClbAccessLogTableIdentifier
}

func (c *ClbAccessLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*ClbAccessLog], error) {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log"),
	}

	return []*table.SourceMetadata[*ClbAccessLog]{
		{
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ClbAccessLog](clbLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ClbAccessLog](clbLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}, nil
}

func (c *ClbAccessLogTable) EnrichRow(row *ClbAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*ClbAccessLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	row.TpIndex = row.Elb

	row.TpSourceIP = &row.ClientIP
	row.TpIps = append(row.TpIps, row.ClientIP)
	if row.BackendIP != nil {
		row.TpDestinationIP = row.BackendIP
		row.TpIps = append(row.TpIps, *row.BackendIP)
	}

	return row, nil
}

func (c *ClbAccessLogTable) GetDescription() string {
	return "AWS CLB access logs capture detailed information about requests processed by a Classic Load Balancer, including client information, backend responses, and SSL details."
}
