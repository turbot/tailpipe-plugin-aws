package tables

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type S3ServerAccessLogTable struct {
	table.TableBase[*S3ServerAccessLogTableConfig]
}

func NewS3ServerAccessLogTable() table.Table {
	return &S3ServerAccessLogTable{}
}

func (c *S3ServerAccessLogTable) Identifier() string {
	return "aws_s3_server_access_log"
}

func (c *S3ServerAccessLogTable) GetRowSchema() any {
	return &AwsS3ServerAccessLog{}
}

func (c *S3ServerAccessLogTable) GetConfigSchema() parse.Config {
	return &S3ServerAccessLogTableConfig{}
}

func (c *S3ServerAccessLogTable) Init(ctx context.Context, tableConfigData *parse.Data, collectionStateJSON json.RawMessage, sourceConfigData *parse.Data) error {
	// call base init
	if err := c.TableBase.Init(ctx, tableConfigData, collectionStateJSON, sourceConfigData); err != nil {
		return err
	}
	// TODO switch on source
	// TODO KAI make sure tables add NewCloudwatchMapper if needed
	// NOTE: add the cloudwatch mapper to ensure rows are in correct format
	//s.AddMappers(artifact_mapper.NewCloudwatchMapper())

	// if the source is an artifact source, we need a mapper
	c.Mapper = mappers.NewS3ServerAccessLogMapper()
	return nil
}

func (c *S3ServerAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *S3ServerAccessLogTable) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// short-circuit for unexpected row type
	rawRecord, ok := row.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid row type: %T, expected map[string]string", row)
	}

	// TODO: #validate ensure we have a timestamp field

	// Build record and add any source enrichment fields
	var record AwsS3ServerAccessLog
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	for key, value := range rawRecord {
		switch key {
		case "bucket_owner":
			record.BucketOwner = value
		case "bucket":
			record.Bucket = value
		case "timestamp":
			ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", value)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestamp: %w", err)
			}
			record.Timestamp = ts
			record.TpTimestamp = helpers.UnixMillis(ts.UnixNano() / int64(time.Millisecond))
		case "remote_ip":
			record.RemoteIP = value
			record.TpSourceIP = &value
			record.TpIps = append(record.TpIps, value)
		case "requester":
			record.Requester = value
		case "request_id":
			record.RequestID = value
		case "operation":
			record.Operation = value
		case "key":
			if value != "-" {
				record.Key = &value
			}
		case "request_uri":
			if value != "-" {
				record.RequestURI = &value
			}
		case "http_status":
			if value != "-" {
				hs, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing http_status: %w", err)
				}
				record.HTTPStatus = &hs
			}
		case "error_code":
			if value != "-" {
				record.ErrorCode = &value
			}
		case "bytes_sent":
			if value != "-" {
				bs, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing bytes_sent: %w", err)
				}
				record.BytesSent = &bs
			}
		case "object_size":
			if value != "-" {
				os, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing object_size: %w", err)
				}
				record.ObjectSize = &os
			}
		case "total_time":
			if value != "-" {
				tt, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing total_time: %w", err)
				}
				record.TotalTime = &tt
			}
		case "turn_around_time":
			if value != "-" {
				tat, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing turn_around_time: %w", err)
				}
				record.TurnAroundTime = &tat
			}
		case "referer":
			if value != "-" {
				record.Referer = &value
			}
		case "user_agent":
			if value != "-" {
				record.UserAgent = &value
			}
		case "version_id":
			if value != "-" {
				record.VersionID = &value
			}
		case "host_id":
			if value != "-" {
				record.HostID = &value
			}
		case "signature_version":
			if value != "-" {
				record.SignatureVersion = &value
			}
		case "cipher_suite":
			if value != "-" {
				record.CipherSuite = &value
			}
		case "authentication_type":
			if value != "-" {
				record.AuthenticationType = &value
			}
		case "host_header":
			if value != "-" {
				record.HostHeader = &value
			}
		case "tls_version":
			if value != "-" {
				record.TLSVersion = &value
			}
		case "access_point_arn":
			if value != "-" {
				record.AccessPointArn = &value
			}
		case "acl_required":
			if value != "-" {
				b := true
				record.AclRequired = &b
			} else {
				b := false
				record.AclRequired = &b
			}
		}
	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	record.TpSourceType = "aws.s3_server_access_log"

	// Hive Fields
	record.TpPartition = c.Identifier()
	if record.TpIndex == "" {
		record.TpIndex = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	// convert to date in format yy-mm-dd
	record.TpDate = record.Timestamp.Format("2006-01-02")

	return record, nil
}
