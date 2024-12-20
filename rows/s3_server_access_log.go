package rows

import (
	"fmt"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type S3ServerAccessLog struct {
	schema.CommonFields

	BucketOwner        string    `json:"bucket_owner"`
	Bucket             string    `json:"bucket"`
	Timestamp          time.Time `json:"timestamp"`
	RemoteIP           string    `json:"remote_ip"`
	Requester          string    `json:"requester"`
	RequestID          string    `json:"request_id"`
	Operation          string    `json:"operation"`
	Key                *string   `json:"key,omitempty"`
	RequestURI         *string   `json:"request_uri,omitempty"`
	HTTPStatus         *int      `json:"http_status,omitempty"`
	ErrorCode          *string   `json:"error_code,omitempty"`
	BytesSent          *int64    `json:"bytes_sent,omitempty"`
	ObjectSize         *int64    `json:"object_size,omitempty"`
	TotalTime          *int      `json:"total_time,omitempty"`
	TurnAroundTime     *int      `json:"turn_around_time,omitempty"`
	Referer            *string   `json:"referer,omitempty"`
	UserAgent          *string   `json:"user_agent,omitempty"`
	VersionID          *string   `json:"version_id,omitempty"`
	HostID             *string   `json:"host_id,omitempty"`
	SignatureVersion   *string   `json:"signature_version,omitempty"`
	CipherSuite        *string   `json:"cipher_suite,omitempty"`
	AuthenticationType *string   `json:"authentication_type,omitempty"`
	HostHeader         *string   `json:"host_header,omitempty"`
	TLSVersion         *string   `json:"tls_version,omitempty"`
	AccessPointArn     *string   `json:"access_point_arn,omitempty"`
	AclRequired        *bool     `json:"acl_required,omitempty"`
}

func NewS3ServerAccessLog() *S3ServerAccessLog {
	return &S3ServerAccessLog{}
}

func (l *S3ServerAccessLog) InitialiseFromMap(m map[string]string) error {
	for key, value := range m {
		switch key {
		case "bucket_owner":
			l.BucketOwner = value
		case "bucket":
			l.Bucket = value
		case "timestamp":
			ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			l.Timestamp = ts
		case "remote_ip":
			l.RemoteIP = value
		case "requester":
			l.Requester = value
		case "request_id":
			l.RequestID = value
		case "operation":
			l.Operation = value
		case "key":
			if value != "-" {
				l.Key = &value
			}
		case "request_uri":
			if value != "-" {
				l.RequestURI = &value
			}
		case "http_status":
			if value != "-" {
				hs, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing http_status: %w", err)
				}
				l.HTTPStatus = &hs
			}
		case "error_code":
			if value != "-" {
				l.ErrorCode = &value
			}
		case "bytes_sent":
			if value != "-" {
				bs, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing bytes_sent: %w", err)
				}
				l.BytesSent = &bs
			}
		case "object_size":
			if value != "-" {
				os, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing object_size: %w", err)
				}
				l.ObjectSize = &os
			}
		case "total_time":
			if value != "-" {
				tt, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing total_time: %w", err)
				}
				l.TotalTime = &tt
			}
		case "turn_around_time":
			if value != "-" {
				tat, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing turn_around_time: %w", err)
				}
				l.TurnAroundTime = &tat
			}
		case "referer":
			if value != "-" {
				l.Referer = &value
			}
		case "user_agent":
			if value != "-" {
				l.UserAgent = &value
			}
		case "version_id":
			if value != "-" {
				l.VersionID = &value
			}
		case "host_id":
			if value != "-" {
				l.HostID = &value
			}
		case "signature_version":
			if value != "-" {
				l.SignatureVersion = &value
			}
		case "cipher_suite":
			if value != "-" {
				l.CipherSuite = &value
			}
		case "authentication_type":
			if value != "-" {
				l.AuthenticationType = &value
			}
		case "host_header":
			if value != "-" {
				l.HostHeader = &value
			}
		case "tls_version":
			if value != "-" {
				l.TLSVersion = &value
			}
		case "access_point_arn":
			if value != "-" {
				l.AccessPointArn = &value
			}
		case "acl_required":
			if value != "-" {
				b := true
				l.AclRequired = &b
			} else {
				b := false
				l.AclRequired = &b
			}
		}
	}
	return nil
}
