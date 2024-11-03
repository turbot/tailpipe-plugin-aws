package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type AwsS3ServerAccessLog struct {
	enrichment.CommonFields

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
