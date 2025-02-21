package vpc_flow_log

import (
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/formats"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const VpcFlowLogTableIdentifier = "aws_vpc_flow_log"

// VpcFlowLogTable - table for AWS VPC Flow Logs
type VpcFlowLogTable struct {
	table.CustomTableImpl
}

func (c *VpcFlowLogTable) Identifier() string {
	return VpcFlowLogTableIdentifier
}

func (c *VpcFlowLogTable) GetFormat() parse.Config {
	return &formats.Regex{
		Layout: `^(?P<version>\d+) (?P<account_id>\d+) (?P<interface_id>\S+) (?P<srcaddr>\S+) (?P<dstaddr>\S+) (?P<srcport>\d+|-) (?P<dstport>\d+|-) (?P<protocol>\d+) (?P<packets>\d+|-) (?P<bytes>\d+|-) (?P<start>\d+) (?P<end>\d+) (?P<action>\S+) (?P<log_status>\S+)`,
	}
}

func (c *VpcFlowLogTable) GetTableDefinition() *schema.TableSchema {
	return &schema.TableSchema{
		Name: VpcFlowLogTableIdentifier,
		Columns: []*schema.ColumnSchema{
			{
				ColumnName: "tp_timestamp",
				SourceName: "start",
			},
			{
				ColumnName: "tp_index",
				SourceName: "account_id",
			},
			{
				ColumnName: "tp_source_ip",
				SourceName: "srcaddr",
			},
			{
				ColumnName: "tp_destination_ip",
				SourceName: "dstaddr",
			},
			// v1/v2 fields
			{
				ColumnName:  "version",
				Description: "VPC Flow Logs version. If you use the default format, the version is 2. If you use a custom format, the version is the highest version among the specified fields.",
				Type:        "INTEGER",
			},
			{
				ColumnName:  "account_id",
				Description: "The AWS account ID of the owner of the source network interface for which traffic is recorded. If the network interface is created by an AWS service, for example when creating a VPC endpoint or Network Load Balancer, the record might display unknown for this field.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "interface_id",
				Description: "The ID of the network interface for which the traffic is recorded.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "src_addr",
				SourceName:  "srcaddr",
				Description: "For incoming traffic, this is the IP address of the source of traffic. For outgoing traffic, this is the private IPv4 address or the IPv6 address of the network interface sending the traffic.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "dst_addr",
				SourceName:  "dstaddr",
				Description: "The destination address for outgoing traffic, or the IPv4 or IPv6 address of the network interface for incoming traffic on the network interface. The IPv4 address of the network interface is always its private IPv4 address.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "src_port",
				SourceName:  "srcport",
				Description: "The source port of the traffic.",
				Type:        "INTEGER",
			},
			{
				ColumnName:  "dst_port",
				SourceName:  "dstport",
				Description: "The destination port of the traffic.",
				Type:        "INTEGER",
			},
			{
				ColumnName:  "protocol",
				Description: "The IANA protocol number of the traffic.",
				Type:        "INTEGER",
			},
			{
				ColumnName:  "packets",
				Description: "The number of packets transferred during the flow.",
				Type:        "BIGINT",
			},
			{
				ColumnName:  "bytes",
				Description: "The number of bytes transferred during the flow.",
				Type:        "BIGINT",
			},
			{
				ColumnName:  "action",
				Description: "The action that is associated with the traffic (ACCEPT/REJECT).",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "log_status",
				Description: "The logging status of the flow log (OK/NODATA/SKIPDATA).",
				Type:        "VARCHAR",
			},
			// v3 fields
			{
				ColumnName:  "vpc_id",
				Description: "The ID of the VPC that contains the network interface for which the traffic is recorded.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "subnet_id",
				Description: "The ID of the subnet that contains the network interface for which the traffic is recorded.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "instance_id",
				Description: "The ID of the instance that's associated with network interface for which the traffic is recorded, if the instance not a requester-managed network interface; for example, the network interface for a NAT gateway.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "tcp_flags",
				Description: "The bitmask value for the following TCP flags: N/A (0) FIN (1), SYN (2), RST (4), SYN-ACK (18).",
				Type:        "INTEGER",
			},
			{
				ColumnName:  "type",
				Description: "The type of traffic (IPv4/IPv6/EFA).",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "pkt_src_addr",
				SourceName:  "pkt_srcaddr",
				Description: "The packet-level (original) source IP address of the traffic.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "pkt_dst_addr",
				SourceName:  "pkt_dstaddr",
				Description: "The packet-level (original) destination IP address for the traffic.",
				Type:        "VARCHAR",
			},
			// v4 fields
			{
				ColumnName:  "region",
				Description: "The Region that contains the network interface for which traffic is recorded.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "az_id",
				Description: "The ID of the Availability Zone that contains the network interface for which traffic is recorded.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "sublocation_type",
				Description: "The type of sublocation that's returned in the sublocation-id field (wavelength/outpost/localzone).",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "sublocation_id",
				Description: "The ID of the sublocation that contains the network interface for which traffic is recorded.",
				Type:        "VARCHAR",
			},
			// v5 fields
			{
				ColumnName:  "pkt_src_aws_service",
				Description: "The name of the subset of IP address ranges for the pkt-srcaddr field, if the source IP address is for an AWS service.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "pkt_dst_aws_service",
				Description: "The name of the subset of IP address ranges for the pkt-dstaddr field, if the destination IP address is for an AWS service.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "flow_direction",
				Description: "The direction of the flow with respect to the interface where traffic is captured (ingress/egress).",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "traffic_path",
				Description: "The path that egress traffic takes to the destination.",
				Type:        "INTEGER",
			},
			// v7 fields
			{
				ColumnName:  "ecs_cluster_arn",
				Description: "AWS Resource Name (ARN) of the ECS cluster if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_cluster_name",
				Description: "Name of the ECS cluster if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_container_instance_arn",
				Description: "ARN of the ECS container instance if the traffic is from a running ECS task on an EC2 instance.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_container_instance_id",
				Description: "ID of the ECS container instance if the traffic is from a running ECS task on an EC2 instance.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_container_id",
				Description: "Docker runtime ID of the container if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_second_container_id",
				Description: "Docker runtime ID of the container if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_service_name",
				Description: "Name of the ECS service if the traffic is from a running ECS task and the ECS task is started by an ECS service.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_task_definition_arn",
				Description: "ARN of the ECS task definition if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_task_arn",
				Description: "ARN of the ECS task if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			{
				ColumnName:  "ecs_task_id",
				Description: "ID of the ECS task if the traffic is from a running ECS task.",
				Type:        "VARCHAR",
			},
			// v8 fields
			{
				ColumnName:  "reject_reason",
				Description: "Reason why traffic was rejected (BPA).",
				Type:        "VARCHAR",
			},
		},
		NullValue: "-", // default null value
	}
}

func (c *VpcFlowLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*table.DynamicRow], error) {
	// ask our CustomTableImpl for the mapper
	mapper, err := c.GetMapper()
	if err != nil {
		return nil, err
	}

	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.gz"),
	}

	return []*table.SourceMetadata[*table.DynamicRow]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
				artifact_source.WithSkipHeaderRow(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
				artifact_source.WithSkipHeaderRow(),
			},
		},
	}, nil
}
