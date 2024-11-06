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

	AccountId          *string       `json:"accountId"`
	Arn                *string       `json:"arn"`
	CreatedAt          time.Time    `json:"createdAt"`
	Description        *string       `json:"description"`
	Id                 *string       `json:"id"`
	Partition          *string       `json:"partition"`
	Region             *string       `json:"region"`
	AccessKeyId        *string       `json:"accessKeyId"`
	PrincipalId        *string       `json:"principalId"`
	UserName           *string       `json:"userName"`
	UserType           *string       `json:"userType"`
	AvailabilityZone   *string       `json:"availabilityZone"`
	InstanceArn        *string       `json:"instanceArn"`
	ImageDescription   *string       `json:"imageDescription"`
	ImageId            *string       `json:"imageId"`
	InstanceId         *string       `json:"instanceId"`
	InstanceState      *string       `json:"instanceState"`
	InstanceType       *string       `json:"instanceType"`
	OutpostArn         *string       `json:"outpostArn"`
	LaunchTime         *string       `json:"launchTime"`
	Ipv6Addresses      []string     `json:"ipv6Addresses"`
	NetworkInterfaceId string       `json:"networkInterfaceId"`
	PrivateDnsName     string       `json:"privateDnsName"`
	PrivateIpAddress   string       `json:"privateIpAddress"`
	PublicDnsName      string       `json:"publicDnsName"`
	PublicIp           string       `json:"publicIp"`
	GroupId            string       `json:"groupId"`
	GroupName          string       `json:"groupName"`
	SubnetId           string       `json:"subnetId"`
	VpcId              string       `json:"vpcId"`
	Platform           *string      `json:"platform,omitempty"`
	Code               string       `json:"code"`
	ProductType        string       `json:"productType"`
	ResourceType       *string       `json:"resourceType"`
	SchemaVersion      *string       `json:"schemaVersion"`
	ActionType         *string       `json:"actionType"`
	Api                *string       `json:"api"`
	CallerType         *string       `json:"callerType"`
	ErrorCode          *string       `json:"errorCode"`
	IpAddressV4        *string       `json:"ipAddressV4"`
	IpAddressV6        *string       `json:"ipAddressV6,omitempty"`
	Archived           *bool         `json:"archived"`
	Count              *int32          `json:"count"`
	DetectorId         *string       `json:"detectorId"`
	EventFirstSeen     *string       `json:"eventFirstSeen"`
	EventLastSeen      *string       `json:"eventLastSeen"`
	ResourceRole       *string       `json:"resourceRole"`
	ServiceName        *string       `json:"serviceName"`
	Severity           *float64      `json:"severity"`
	Title              *string       `json:"title"`
	Type               *string       `json:"type"`
	UpdatedAt          *string       `json:"updatedAt"`
}
