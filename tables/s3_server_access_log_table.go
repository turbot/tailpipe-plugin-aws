package tables

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type S3ServerAccessLogTable struct {
	table.TableImpl[map[string]string, *S3ServerAccessLogTableConfig, *config.AwsConnection]
}

func NewS3ServerAccessLogTable() table.Table {
	return &S3ServerAccessLogTable{}
}

func (c *S3ServerAccessLogTable) Identifier() string {
	return "aws_s3_server_access_log"
}

func (c *S3ServerAccessLogTable) GetRowSchema() any {
	return &rows.AwsS3ServerAccessLog{}
}

func (c *S3ServerAccessLogTable) GetConfigSchema() parse.Config {
	return &S3ServerAccessLogTableConfig{}
}

func (c *S3ServerAccessLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMappers()
	return nil
}

func (c *S3ServerAccessLogTable) initMappers() {
	// TODO switch on source

	// TODO KAI make sure tables add NewCloudwatchMapper if needed
	// NOTE: add the cloudwatch mapper to ensure rows are in correct format
	//s.AddMappers(artifact_mapper.NewCloudwatchMapper())

	// if the source is an artifact source, we need a mapper
	c.Mapper = mappers.NewS3ServerAccessLogMapper()
}

func (c *S3ServerAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *S3ServerAccessLogTable) EnrichRow(rawRow map[string]string, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// TODO: #validate ensure we have a timestamp field

	// Build row and add any source enrichment fields
	var row rows.AwsS3ServerAccessLog
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	for key, value := range rawRow {
		switch key {
		case "bucket_owner":
			row.BucketOwner = value
		case "bucket":
			row.Bucket = value
		case "timestamp":
			ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", value)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestamp: %w", err)
			}
			row.Timestamp = ts
			row.TpTimestamp = helpers.UnixMillis(ts.UnixNano() / int64(time.Millisecond))
		case "remote_ip":
			row.RemoteIP = value
			row.TpSourceIP = &value
			row.TpIps = append(row.TpIps, value)
		case "requester":
			row.Requester = value
		case "request_id":
			row.RequestID = value
		case "operation":
			row.Operation = value
		case "key":
			if value != "-" {
				row.Key = &value
			}
		case "request_uri":
			if value != "-" {
				row.RequestURI = &value
			}
		case "http_status":
			if value != "-" {
				hs, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing http_status: %w", err)
				}
				row.HTTPStatus = &hs
			}
		case "error_code":
			if value != "-" {
				row.ErrorCode = &value
			}
		case "bytes_sent":
			if value != "-" {
				bs, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing bytes_sent: %w", err)
				}
				row.BytesSent = &bs
			}
		case "object_size":
			if value != "-" {
				os, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing object_size: %w", err)
				}
				row.ObjectSize = &os
			}
		case "total_time":
			if value != "-" {
				tt, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing total_time: %w", err)
				}
				row.TotalTime = &tt
			}
		case "turn_around_time":
			if value != "-" {
				tat, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing turn_around_time: %w", err)
				}
				row.TurnAroundTime = &tat
			}
		case "referer":
			if value != "-" {
				row.Referer = &value
			}
		case "user_agent":
			if value != "-" {
				row.UserAgent = &value
			}
		case "version_id":
			if value != "-" {
				row.VersionID = &value
			}
		case "host_id":
			if value != "-" {
				row.HostID = &value
			}
		case "signature_version":
			if value != "-" {
				row.SignatureVersion = &value
			}
		case "cipher_suite":
			if value != "-" {
				row.CipherSuite = &value
			}
		case "authentication_type":
			if value != "-" {
				row.AuthenticationType = &value
			}
		case "host_header":
			if value != "-" {
				row.HostHeader = &value
			}
		case "tls_version":
			if value != "-" {
				row.TLSVersion = &value
			}
		case "access_point_arn":
			if value != "-" {
				row.AccessPointArn = &value
			}
		case "acl_required":
			if value != "-" {
				b := true
				row.AclRequired = &b
			} else {
				b := false
				row.AclRequired = &b
			}
		}
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
