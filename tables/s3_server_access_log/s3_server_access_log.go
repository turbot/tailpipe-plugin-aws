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

func (c *S3ServerAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"access_point_arn":    "The ARN of the S3 Access Point used for the request, if applicable.",
		"acl_required":        "Indicates if ACLs were required for the request (true/false).",
		"authentication_type": "The authentication method used (e.g., AuthHeader, QueryString).",
		"bucket":              "The name of the S3 bucket where the request was made.",
		"bucket_owner":        "The canonical user ID of the bucket owner.",
		"bytes_sent":          "The number of bytes sent in response to the request.",
		"cipher_suite":        "The cipher suite used for SSL/TLS connections.",
		"error_code":          "The error code returned, if the request resulted in an error.",
		"host_header":         "The host header included in the request.",
		"host_id":             "An identifier assigned to help diagnose request failures.",
		"http_status":         "The HTTP status code returned for the request.",
		"key":                 "The object key (name) if applicable to the request.",
		"object_size":         "The size of the requested object in bytes.",
		"operation":           "The type of operation performed on the S3 object (e.g., REST.GET.OBJECT).",
		"referer":             "The referer header from the client request, if present.",
		"remote_ip":           "The IP address of the client that made the request.",
		"request_id":          "A unique identifier assigned to the request by AWS.",
		"request_uri":         "The full request URI used in the operation.",
		"requester":           "The canonical user ID or IAM role of the entity making the request.",
		"signature_version":   "The signature version used for request authentication (e.g., SigV2, SigV4).",
		"timestamp":           "The date and time when the request was received by S3.",
		"tls_version":         "The TLS protocol version used for the request.",
		"total_time":          "The total time taken for the request from start to finish (in milliseconds).",
		"turn_around_time":    "The time taken by S3 to process the request (in milliseconds).",
		"user_agent":          "The User-Agent string from the client making the request.",
		"version_id":          "The version ID of the object, if versioning is enabled.",

		// Tailpipe-specific metadata fields
		"tp_ips":       "All IP addresses associated with the request, including the remote IP.",
		"tp_usernames": "Canonical user IDs or role ARNs associated with the request.",
	}
}
