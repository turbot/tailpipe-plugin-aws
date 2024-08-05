package aws_collection

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/hcl"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"strconv"
	"time"
)

type S3ServerAccessLogCollection struct {
	collection.CollectionBase[*S3ServerAccessLogCollectionConfig]
}

func NewS3ServerAccessLogCollection() collection.Collection {
	return &S3ServerAccessLogCollection{}
}

func (c *S3ServerAccessLogCollection) Identifier() string {
	return "aws_s3_server_access_log"
}

func (c *S3ServerAccessLogCollection) GetRowSchema() any {
	return &aws_types.AwsS3ServerAccessLog{}
}

func (c *S3ServerAccessLogCollection) GetConfigSchema() hcl.Config {
	return &S3ServerAccessLogCollectionConfig{}
}

func (c *S3ServerAccessLogCollection) GetSourceOptions() []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
		artifact_source.WithMapper(aws_source.NewS3ServerAccessLogMapper()),
	}
}

func (c *S3ServerAccessLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// short-circuit for unexpected row type
	rawRecord, ok := row.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid row type: %T, expected map[string]string", row)
	}

	// TODO: #validate ensure we have a timestamp field

	// Build record and add any source enrichment fields
	var record aws_types.AwsS3ServerAccessLog
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
	record.TpCollection = c.Identifier()
	if record.TpConnection == "" {
		record.TpConnection = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	record.TpYear = int32(record.Timestamp.Year())
	record.TpMonth = int32(record.Timestamp.Month())
	record.TpDay = int32(record.Timestamp.Day())

	return record, nil
}
