package s3_server_access_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"

	"github.com/rs/xid"

	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const S3ServerAccessLogTableIdentifier = "aws_s3_server_access_log"
const s3ServerAccessLogFormat = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn $acl_required`
const s3ServerAccessLogFormatReduced = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn`

type S3ServerAccessLogTable struct{}

func (c *S3ServerAccessLogTable) Identifier() string {
	return S3ServerAccessLogTableIdentifier
}

func (c *S3ServerAccessLogTable) GetSourceMetadata() []*table.SourceMetadata[*S3ServerAccessLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		// FileLayout: utils.ToStringPointer("AWSLogs/%{DATA:bucket_name}/%{YEAR:year}-%{MONTHNUM:month}-%{MONTHDAY:day}-%{HOUR:hour}-%{MINUTE:minute}-%{SECOND:second}-%{DATA:suffix}.gz"),
		FileLayout: utils.ToStringPointer("%{YEAR:year}-%{MONTHNUM:month}-%{MONTHDAY:day}-%{HOUR:hour}-%{MINUTE:minute}-%{SECOND:second}-%{DATA:suffix}"),
	}

	return []*table.SourceMetadata[*S3ServerAccessLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*S3ServerAccessLog](s3ServerAccessLogFormat, s3ServerAccessLogFormatReduced),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*S3ServerAccessLog](s3ServerAccessLogFormat, s3ServerAccessLogFormatReduced),
			Options:    []row_source.RowSourceOption{artifact_source.WithRowPerLine()},
		},
	}
}

func (c *S3ServerAccessLogTable) EnrichRow(row *S3ServerAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*S3ServerAccessLog, error) {
	// TODO: #validate ensure we have a timestamp field

	// add any source enrichment fields
	row.CommonFields = sourceEnrichmentFields.CommonFields

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
