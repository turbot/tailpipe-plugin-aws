package tables

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const s3ServerAccessLogFormat = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn $acl_required`
const s3ServerAccessLogFormatReduced = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn`

type S3ServerAccessLogTable struct {
	table.TableImpl[*rows.S3ServerAccessLog, *S3ServerAccessLogTableConfig, *config.AwsConnection]
}

func NewS3ServerAccessLogTable() table.Table {
	return &S3ServerAccessLogTable{}
}

func (c *S3ServerAccessLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMapper()
	return nil
}

func (c *S3ServerAccessLogTable) initMapper() {
	// TODO switch on source
	c.Mapper = table.NewDelimitedLineMapper(rows.NewS3ServerAccessLog, s3ServerAccessLogFormat, s3ServerAccessLogFormatReduced)
}

func (c *S3ServerAccessLogTable) Identifier() string {
	return "aws_s3_server_access_log"
}

func (c *S3ServerAccessLogTable) GetRowSchema() any {
	return rows.NewS3ServerAccessLog()
}

func (c *S3ServerAccessLogTable) GetConfigSchema() parse.Config {
	return &S3ServerAccessLogTableConfig{}
}

func (c *S3ServerAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *S3ServerAccessLogTable) EnrichRow(row *rows.S3ServerAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// TODO: #validate ensure we have a timestamp field

	// add any source enrichment fields
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	row.TpSourceType = "aws.s3_server_access_log"

	// Hive Fields
	row.TpPartition = c.Identifier()
	if row.TpIndex == "" {
		row.TpIndex = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	// convert to date in format yy-mm-dd
	row.TpDate = row.Timestamp.Format("2006-01-02")

	return row, nil
}
