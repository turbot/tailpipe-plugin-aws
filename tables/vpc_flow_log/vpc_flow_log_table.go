package vpc_flow_log

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch_log_group"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/error_types"
	"github.com/turbot/tailpipe-plugin-sdk/formats"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const VpcFlowLogTableIdentifier = "aws_vpc_flow_log"
const VpcFlowLogTableNilValue = "-"
const VpcFlowLogTableSkippedData = "SKIPDATA"
const VpcFlowLogTableNoData = "NODATA"

// VpcFlowLogTable - table for VPC Flow Logs
type VpcFlowLogTable struct {
	table.CustomTableImpl
}

func (c *VpcFlowLogTable) Identifier() string {
	return VpcFlowLogTableIdentifier
}

func (c *VpcFlowLogTable) GetDefaultFormat() formats.Format {
	return defaultVPCFlowLogTableFormat
}

func (c *VpcFlowLogTable) GetTableDefinition() *schema.TableSchema {
	return &schema.TableSchema{
		Name: VpcFlowLogTableIdentifier,
		Columns: []*schema.ColumnSchema{
			// version 2 (default) fields
			{
				ColumnName:  "version",
				Description: "The VPC Flow Logs version. The version depends on the fields included in the log.",
				Type:        "integer",
			},
			{
				ColumnName:  "account_id",
				Description: "The AWS account ID of the network interface owner.",
				Type:        "varchar",
			},
			{
				ColumnName:  "interface_id",
				Description: "The ID of the network interface for which traffic is recorded.",
				Type:        "varchar",
			},
			{
				ColumnName:  "src_addr",
				Description: "The source IP address of the traffic. For outgoing traffic, this is the private IP of the network interface.",
				Type:        "varchar",
			},
			{
				ColumnName:  "dst_addr",
				Description: "The destination IP address of the traffic. For incoming traffic, this is the private IP of the network interface.",
				Type:        "varchar",
			},
			{
				ColumnName:  "src_port",
				Description: "The source port of the traffic.",
				Type:        "integer",
			},
			{
				ColumnName:  "dst_port",
				Description: "The destination port of the traffic.",
				Type:        "integer",
			},
			{
				ColumnName:  "protocol",
				Description: "The IANA protocol number of the traffic (e.g., 6 for TCP, 17 for UDP).",
				Type:        "integer",
			},
			{
				ColumnName:  "packets",
				Description: "The number of packets transferred during the flow.",
				Type:        "bigint",
			},
			{
				ColumnName:  "bytes",
				Description: "The number of bytes transferred during the flow.",
				Type:        "bigint",
			},
			{
				ColumnName:  "start_time",
				Description: "The time, in Unix seconds, when the first packet of the flow was received within the aggregation interval.",
				Type:        "timestamp",
			},
			{
				ColumnName:  "end_time",
				Description: "The time, in Unix seconds, when the last packet of the flow was received within the aggregation interval.",
				Type:        "timestamp",
			},
			{
				ColumnName:  "action",
				Description: "The action associated with the traffic: ACCEPT or REJECT.",
				Type:        "varchar",
			},
			{
				ColumnName:  "log_status",
				Description: "The logging status of the flow log: OK, NODATA, or SKIPDATA.",
				Type:        "varchar",
			},
			// version 3 fields
			{
				ColumnName:  "vpc_id",
				Description: "The ID of the VPC that contains the network interface for which traffic is recorded.",
				Type:        "varchar",
			},
			{
				ColumnName:  "subnet_id",
				Description: "The ID of the subnet that contains the network interface for which traffic is recorded.",
				Type:        "varchar",
			},
			{
				ColumnName:  "instance_id",
				Description: "The ID of the EC2 instance associated with the network interface, if applicable.",
				Type:        "varchar",
			},
			{
				ColumnName:  "tcp_flags",
				Description: "The bitmask value for TCP flags recorded during the flow.",
				Type:        "integer",
			},
			{
				ColumnName:  "type",
				Description: "The type of traffic (IPv4, IPv6, or EFA).",
				Type:        "varchar",
			},
			{
				ColumnName:  "pkt_src_addr",
				Description: "The packet-level (original) source IP address of the traffic.",
				Type:        "varchar",
			},
			{
				ColumnName:  "pkt_dst_addr",
				Description: "The packet-level (original) destination IP address of the traffic.",
				Type:        "varchar",
			},
			// version 4 fields
			{
				ColumnName:  "region",
				Description: "The AWS region containing the network interface for which traffic is recorded.",
				Type:        "varchar",
			},
			{
				ColumnName:  "az_id",
				Description: "The Availability Zone ID that contains the network interface for which traffic is recorded.",
				Type:        "varchar",
			},
			{
				ColumnName:  "sublocation_type",
				Description: "The type of sublocation that contains the network interface (e.g., wavelength, outpost, localzone).",
				Type:        "varchar",
			},
			{
				ColumnName:  "sublocation_id",
				Description: "The ID of the sublocation that contains the network interface.",
				Type:        "varchar",
			},
			// version 5 fields
			{
				ColumnName:  "pkt_src_aws_service",
				Description: "The AWS service associated with the packet source IP, if applicable.",
				Type:        "varchar",
			},
			{
				ColumnName:  "pkt_dst_aws_service",
				Description: "The AWS service associated with the packet destination IP, if applicable.",
				Type:        "varchar",
			},
			{
				ColumnName:  "flow_direction",
				Description: "The direction of the flow with respect to the network interface: ingress or egress.",
				Type:        "varchar",
			},
			{
				ColumnName:  "traffic_path",
				Description: "The path that egress traffic takes to the destination.",
				Type:        "integer",
			},
			// version 7 fields
			{
				ColumnName:  "ecs_cluster_arn",
				Description: "The ARN of the ECS cluster if traffic originates from an ECS task.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_cluster_name",
				Description: "The name of the ECS cluster if traffic originates from an ECS task.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_container_instance_arn",
				Description: "The ARN of the ECS container instance if traffic is from a running ECS task on an EC2 instance.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_container_instance_id",
				Description: "The ID of the ECS container instance if traffic is from a running ECS task on an EC2 instance.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_container_id",
				Description: "The Docker runtime ID of the first container in an ECS task.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_second_container_id",
				Description: "The Docker runtime ID of the second container in an ECS task, if applicable.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_service_name",
				Description: "The name of the ECS service if traffic is from an ECS task started by a service.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_task_arn",
				Description: "The ARN of the ECS task if traffic is from a running ECS task.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_task_definition_arn",
				Description: "The ARN of the ECS task definition if traffic is from a running ECS task.",
				Type:        "varchar",
			},
			{
				ColumnName:  "ecs_task_id",
				Description: "The ID of the ECS task if traffic is from a running ECS task.",
				Type:        "varchar",
			},
			// version 8 fields
			{
				ColumnName:  "reject_reason",
				Description: "The reason why traffic was rejected (e.g., BPA).",
				Type:        "varchar",
			},
		},
		NullValue:   VpcFlowLogTableNilValue,
		Description: c.GetDescription(),
	}
}

func (c *VpcFlowLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*types.DynamicRow], error) {

	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz"),
	}

	// ask our CustomTableImpl for the mapper
	mapper, err := c.Format.GetMapper()
	if err != nil {
		return nil, err
	}

	return []*table.SourceMetadata[*types.DynamicRow]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &VPCFlowLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewVPCFlowLogExtractor(c.Format)),
			},
		},
		{
			// CloudWatch source
			SourceName: cloudwatch_log_group.AwsCloudwatchSourceIdentifier,
			Mapper:     mapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}, nil
}

// EnrichRow implements table.Table
func (c *VpcFlowLogTable) EnrichRow(row *types.DynamicRow, sourceEnrichmentFields schema.SourceEnrichment) (*types.DynamicRow, error) {
	var invalidFields []string

	row.OutputColumns[constants.TpTable] = VpcFlowLogTableIdentifier

	if startTime, ok := row.GetSourceValue("start_time"); ok && (startTime != VpcFlowLogTableSkippedData || startTime != VpcFlowLogTableNoData || startTime != VpcFlowLogTableNilValue) {
		t, err := helpers.ParseTime(startTime)
		if err != nil {
			invalidFields = append(invalidFields, "start_time")
		} else {
			row.OutputColumns[constants.TpTimestamp] = t
		}
	} else if endTime, ok := row.GetSourceValue("end_time"); ok && (endTime != VpcFlowLogTableSkippedData || endTime != VpcFlowLogTableNoData || endTime != VpcFlowLogTableNilValue) {
		t, err := helpers.ParseTime(endTime)
		if err != nil {
			invalidFields = append(invalidFields, "end_time")
		} else if row.OutputColumns[constants.TpTimestamp] == nil {
			row.OutputColumns[constants.TpTimestamp] = t
		}
	}

	if len(invalidFields) > 0 {
		return nil, error_types.NewRowErrorWithFields([]string{}, invalidFields)
	}

	// tp_ips
	var ips []string
	if srcAddr, ok := row.GetSourceValue("src_addr"); ok && (srcAddr != VpcFlowLogTableSkippedData || srcAddr != VpcFlowLogTableNoData || srcAddr != VpcFlowLogTableNilValue) {
		ips = append(ips, srcAddr)
		row.OutputColumns[constants.TpSourceIP] = srcAddr
	}
	if pktSrcAddr, ok := row.GetSourceValue("pkt_src_addr"); ok && (pktSrcAddr != VpcFlowLogTableSkippedData || pktSrcAddr != VpcFlowLogTableNoData || pktSrcAddr != VpcFlowLogTableNilValue) {
		ips = append(ips, pktSrcAddr)
	}
	if dstAddr, ok := row.GetSourceValue("dst_addr"); ok && (dstAddr != VpcFlowLogTableSkippedData || dstAddr != VpcFlowLogTableNoData || dstAddr != VpcFlowLogTableNilValue) {
		ips = append(ips, dstAddr)
		row.OutputColumns[constants.TpDestinationIP] = dstAddr
	}
	if pktDstAddr, ok := row.GetSourceValue("pkt_dst_addr"); ok && (pktDstAddr != VpcFlowLogTableSkippedData || pktDstAddr != VpcFlowLogTableNoData || pktDstAddr != VpcFlowLogTableNilValue) {
		ips = append(ips, pktDstAddr)
	}
	if len(ips) > 0 {
		row.OutputColumns[constants.TpIps] = ips
	}

	// tp_index
	for _, key := range []string{"interface_id", "subnet_id", "vpc_id"} {
		if val, ok := row.GetSourceValue(key); ok && (val != VpcFlowLogTableSkippedData || val != VpcFlowLogTableNoData || val != VpcFlowLogTableNilValue) {
			row.OutputColumns[constants.TpIndex] = val
			break
		}
	}
	if row.OutputColumns[constants.TpIndex] == nil {
		row.OutputColumns[constants.TpIndex] = "default"
	}

	// tp_akas
	var akas []string
	if ecsClusterArn, ok := row.GetSourceValue("ecs_cluster_arn"); ok && (ecsClusterArn != VpcFlowLogTableSkippedData || ecsClusterArn != VpcFlowLogTableNoData || ecsClusterArn != VpcFlowLogTableNilValue) {
		akas = append(akas, ecsClusterArn)
	}
	if ecsContainerInstanceArn, ok := row.GetSourceValue("ecs_container_instance_arn"); ok && (ecsContainerInstanceArn != VpcFlowLogTableSkippedData || ecsContainerInstanceArn != VpcFlowLogTableNoData || ecsContainerInstanceArn != VpcFlowLogTableNilValue) {
		akas = append(akas, ecsContainerInstanceArn)
	}
	if ecsTaskArn, ok := row.GetSourceValue("ecs_task_arn"); ok && (ecsTaskArn != VpcFlowLogTableSkippedData || ecsTaskArn != VpcFlowLogTableNoData || ecsTaskArn != VpcFlowLogTableNilValue) {
		akas = append(akas, ecsTaskArn)
	}
	if ecsTaskDefinitionArn, ok := row.GetSourceValue("ecs_task_definition_arn"); ok && (ecsTaskDefinitionArn != VpcFlowLogTableSkippedData || ecsTaskDefinitionArn != VpcFlowLogTableNoData || ecsTaskDefinitionArn != VpcFlowLogTableNilValue) {
		akas = append(akas, ecsTaskDefinitionArn)
	}
	if len(akas) > 0 {
		row.OutputColumns[constants.TpAkas] = akas
	}

	// now call the base class to do the rest of the enrichment
	return c.CustomTableImpl.EnrichRow(row, sourceEnrichmentFields)
}

func (c *VpcFlowLogTable) GetDescription() string {
	return "AWS VPC Flow Logs capture information about IP traffic going to and from network interfaces in your VPC. This table provides detailed network traffic patterns, including source and destination IP addresses, ports, protocols, and traffic volumes, helping teams monitor network flows, troubleshoot connectivity issues, and detect security anomalies."
}
