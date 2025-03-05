package network_firewall_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const NetworkFirewallLogTableIdentifier = "aws_network_firewall_log"

// NetworkFirewallLogTable implements the table interface for AWS Network Firewall logs.
type NetworkFirewallLogTable struct{}

// Identifier returns the unique table identifier.
func (c *NetworkFirewallLogTable) Identifier() string {
	return NetworkFirewallLogTableIdentifier
}

// GetSourceMetadata returns the artifact source configurations for the table.
func (c *NetworkFirewallLogTable) GetSourceMetadata() []*table.SourceMetadata[*NetworkFirewallLog] {
	// Example file layout – adjust as needed.
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/networkfirewall/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log"),
	}

	return []*table.SourceMetadata[*NetworkFirewallLog]{
		{
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			// Use JSON mapping since the logs are structured as JSON.
			Mapper: &NetworkFirewallMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &NetworkFirewallMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

// EnrichRow applies standard tailpipe enrichment to each log record.
func (c *NetworkFirewallLogTable) EnrichRow(row *NetworkFirewallLog, sourceEnrichmentFields schema.SourceEnrichment) (*NetworkFirewallLog, error) {
	// Add common enrichment fields.
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Standard tailpipe fields.
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	// Convert the epoch event_timestamp (in seconds) to time.Time.
	row.TpTimestamp = *row.EventTimestamp
	row.TpDate = row.TpTimestamp.Truncate(24 * time.Hour)

	// Use the firewall name as an index.
	row.TpIndex = row.FirewallName

	// If available, use the source IP from the event details.
	if row.Event.SrcIP != "" {
		row.TpSourceIP = &row.Event.SrcIP
		row.TpIps = append(row.TpIps, row.Event.SrcIP)
	}

	return row, nil
}

// GetDescription returns a human-readable description of the table.
func (c *NetworkFirewallLogTable) GetDescription() string {
	return "AWS Network Firewall logs capture detailed information about firewall events—including alert, flow, and TLS events produced by Suricata and a dedicated TLS engine. This table provides a structured representation of the log data."
}