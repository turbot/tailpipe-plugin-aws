package s3_server_access_log

import (
	"fmt"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type S3ServerAccessLog struct {
	schema.CommonFields

	AccessPointArn     *string   `json:"access_point_arn,omitempty"`
	AclRequired        *bool     `json:"acl_required,omitempty"`
	AuthenticationType *string   `json:"authentication_type,omitempty"`
	Bucket             string    `json:"bucket"`
	BucketOwner        string    `json:"bucket_owner"`
	BytesSent          *int64    `json:"bytes_sent,omitempty"`
	CipherSuite        *string   `json:"cipher_suite,omitempty"`
	ErrorCode          *string   `json:"error_code,omitempty"`
	HTTPStatus         *int      `json:"http_status"`
	HostHeader         *string   `json:"host_header,omitempty"`
	HostID             *string   `json:"host_id,omitempty"`
	Key                *string   `json:"key,omitempty"`
	ObjectSize         *int64    `json:"object_size,omitempty"`
	Operation          string    `json:"operation"`
	Referer            *string   `json:"referer,omitempty"`
	RemoteIP           string    `json:"remote_ip"`
	RequestID          string    `json:"request_id"`
	RequestURI         *string   `json:"request_uri"`
	Requester          string    `json:"requester,omitempty"`
	SignatureVersion   *string   `json:"signature_version,omitempty"`
	TLSVersion         *string   `json:"tls_version,omitempty"`
	Timestamp          time.Time `json:"timestamp"`
	TotalTime          *int      `json:"total_time"`
	TurnAroundTime     *int      `json:"turn_around_time,omitempty"`
	UserAgent          *string   `json:"user_agent,omitempty"`
	VersionID          *string   `json:"version_id,omitempty"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *S3ServerAccessLog) InitialiseFromMap(m map[string]string) error {
	var err error
	for key, value := range m {
		if value == "-" {
			continue
		}
		switch key {
		case "bucket_owner":
			l.BucketOwner = value
		case "bucket":
			l.Bucket = value
		case "timestamp":
			l.Timestamp, err = time.Parse("02/Jan/2006:15:04:05 -0700", value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
		case "remote_ip":
			l.RemoteIP = value
		case "request_id":
			l.RequestID = value
		case "operation":
			l.Operation = value
		case "requester":
			l.Requester = value
		case "key":
			l.Key = &value
		case "request_uri":
			l.RequestURI = &value
		case "error_code":
			l.ErrorCode = &value
		case "referer":
			l.Referer = &value
		case "user_agent":
			l.UserAgent = &value
		case "version_id":
			l.VersionID = &value
		case "host_id":
			l.HostID = &value
		case "signature_version":
			l.SignatureVersion = &value
		case "cipher_suite":
			l.CipherSuite = &value
		case "authentication_type":
			l.AuthenticationType = &value
		case "host_header":
			l.HostHeader = &value
		case "tls_version":
			l.TLSVersion = &value
		case "access_point_arn":
			l.AccessPointArn = &value
		case "http_status":
			hs, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing http_status: %w", err)
			}
			l.HTTPStatus = &hs
		case "bytes_sent":
			bs, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing bytes_sent: %w", err)
			}
			l.BytesSent = &bs
		case "object_size":
			os, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing object_size: %w", err)
			}
			l.ObjectSize = &os
		case "total_time":
			tt, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing total_time: %w", err)
			}
			l.TotalTime = &tt
		case "turn_around_time":
			tat, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing turn_around_time: %w", err)
			}
			l.TurnAroundTime = &tat
		case "acl_required":
			b := true
			l.AclRequired = &b
		}
	}
	return nil
}
