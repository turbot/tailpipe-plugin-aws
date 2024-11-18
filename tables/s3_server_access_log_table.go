package tables

import (
	"time"

	"github.com/rs/xid"

	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const S3ServerAccessLogTableIdentifier = "aws_s3_server_access_log"
const s3ServerAccessLogFormat = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn $acl_required`
const s3ServerAccessLogFormatReduced = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn`

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.S3ServerAccessLog, *S3ServerAccessLogTable]()
}

type S3ServerAccessLogTable struct {
	table.TableImpl[*rows.S3ServerAccessLog, *S3ServerAccessLogTableConfig, *config.AwsConnection]
}

func (c *S3ServerAccessLogTable) Identifier() string {
	return S3ServerAccessLogTableIdentifier
}

func (c *S3ServerAccessLogTable) SupportedSources() []*table.SourceMetadata[*rows.S3ServerAccessLog] {
	return []*table.SourceMetadata[*rows.S3ServerAccessLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: c.initMapper(),
			Options:    []row_source.RowSourceOption{artifact_source.WithRowPerLine()},
		},
	}
}

func (c *S3ServerAccessLogTable) initMapper() func() table.Mapper[*rows.S3ServerAccessLog] {
	f := func() table.Mapper[*rows.S3ServerAccessLog] {
		return table.NewDelimitedLineMapper(rows.NewS3ServerAccessLog, s3ServerAccessLogFormat, s3ServerAccessLogFormatReduced)
	}
	return f
}

func (c *S3ServerAccessLogTable) EnrichRow(row *rows.S3ServerAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.S3ServerAccessLog, error) {
	// TODO: #validate ensure we have a timestamp field

	// add any source enrichment fields
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	row.TpIndex = row.Bucket // TODO: #enrichment this would ideally be the AccountID, how to obtain?

	// IPs
	row.TpSourceIP = &row.RemoteIP
	row.TpIps = append(row.TpIps, row.RemoteIP)

	row.TpUsernames = append(row.TpUsernames, row.Requester, row.BucketOwner)

	return row, nil
}
