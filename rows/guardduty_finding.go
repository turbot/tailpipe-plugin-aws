package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type GuardDutyBatch struct {
	Records []GuardDutyFinding `json:"Records"`
}

type GuardDutyFinding struct {
	enrichment.CommonFields

	AccountId          *string    `json:"account_id"`
	Arn                *string    `json:"arn"`
	CreatedAt          time.Time  `json:"created_at"`
	Description        *string    `json:"description"`
	Id                 *string    `json:"id"`
	Partition          *string    `json:"partition"`
	Region             *string    `json:"region"`
	AccessKeyId        *string    `json:"access_key_id"`
	PrincipalId        *string    `json:"principal_id"`
	UserName           *string    `json:"user_name"`
	UserType           *string    `json:"user_type"`
	AvailabilityZone   *string    `json:"availability_zone"`
	InstanceArn        *string    `json:"instance_arn"`
	ImageDescription   *string    `json:"image_description"`
	ImageId            *string    `json:"image_id"`
	InstanceId         *string    `json:"instance_id"`
	InstanceState      *string    `json:"instance_state"`
	InstanceType       *string    `json:"instance_type"`
	OutpostArn         *string    `json:"outpost_arn"`
	LaunchTime         *time.Time `json:"launch_time"`
	Ipv6Addresses      []string   `json:"ipv6_addresses" parquet:"type=JSON"`
	NetworkInterfaceId string     `json:"network_interface_id"`
	PrivateDnsName     string     `json:"private_dns_name"`
	PrivateIpAddress   string     `json:"private_ip_address"`
	PublicDnsName      string     `json:"public_dns_name"`
	PublicIp           string     `json:"public_ip"`
	GroupId            string     `json:"group_id"`
	GroupName          string     `json:"group_name"`
	SubnetId           string     `json:"subnet_id"`
	VpcId              string     `json:"vpc_id"`
	Platform           *string    `json:"platform,omitempty"`
	Code               string     `json:"code"`
	ProductType        string     `json:"product_type"`
	ResourceType       *string    `json:"resource_type"`
	SchemaVersion      *string    `json:"schema_version"`
	ActionType         *string    `json:"action_type"`
	Api                *string    `json:"api"`
	CallerType         *string    `json:"caller_type"`
	ErrorCode          *string    `json:"error_code"`
	IpAddressV4        *string    `json:"ip_address_v4"`
	IpAddressV6        *string    `json:"ip_address_v6,omitempty"`
	Archived           *bool      `json:"archived"`
	Count              *int32     `json:"count"`
	DetectorId         *string    `json:"detector_id"`
	EventFirstSeen     *string    `json:"event_first_seen"`
	EventLastSeen      *string    `json:"event_last_seen"`
	ResourceRole       *string    `json:"resource_role"`
	ServiceName        *string    `json:"service_name"`
	Severity           *float64   `json:"severity"`
	Title              *string    `json:"title"`
	Type               *string    `json:"type"`
	UpdatedAt          *time.Time `json:"updated_at"`
}
