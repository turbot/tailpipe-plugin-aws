package nlb_access_log

import (
	"strings"
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

const NlbAccessLogTableIdentifier = "aws_nlb_access_log"

const nlbLogFormat = `$type $version $timestamp $elb $listener $client_ip:$client_port $destination_ip:$destination_port $connection_time $tls_handshake_time $received_bytes $sent_bytes $incoming_tls_alert $chosen_cert_arn $chosen_cert_serial $tls_cipher $tls_protocol_version $tls_named_group $domain_name $alpn_fe_protocol $alpn_be_protocol $alpn_client_preference_list $tls_connection_creation_time`

type NlbAccessLogTable struct{}

func (c *NlbAccessLogTable) Identifier() string {
	return NlbAccessLogTableIdentifier
}

func (c *NlbAccessLogTable) GetSourceMetadata() []*table.SourceMetadata[*NlbAccessLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.%{DATA}.log.gz"),
	}

	return []*table.SourceMetadata[*NlbAccessLog]{
		{
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*NlbAccessLog](nlbLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*NlbAccessLog](nlbLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *NlbAccessLogTable) EnrichRow(row *NlbAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*NlbAccessLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	row.TpIndex = strings.TrimPrefix(row.Elb, "net/")

	row.TpSourceIP = &row.ClientIP
	row.TpIps = append(row.TpIps, row.ClientIP)
	if row.DestinationIP != "" {
		row.TpDestinationIP = &row.DestinationIP
		row.TpIps = append(row.TpIps, row.DestinationIP)
	}

	if row.DomainName != "" {
		row.TpDomains = append(row.TpDomains, row.DomainName)
	}

	if row.ChosenCertArn != "" {
		row.TpAkas = append(row.TpAkas, row.ChosenCertArn)
	}

	return row, nil
}

func (c *NlbAccessLogTable) GetDescription() string {
	return "AWS NLB access logs capture detailed information about the connections that pass through a Network Load Balancer. This table provides a structured representation of the log data."
}
