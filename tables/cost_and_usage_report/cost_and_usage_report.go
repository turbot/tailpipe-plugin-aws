package cost_and_usage_report

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// The following structure is defined based on the legacy and CUR 2.0 Schema structure.
// Legacy Schema: https://docs.aws.amazon.com/cur/latest/userguide/data-dictionary.html
// CUR 2.0 Schema: https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cur2.html
type CostAndUsageLog struct {
	schema.CommonFields

	BillBillingEntity                                        *string                 `json:"bill_billing_entity,omitempty" parquet:"name=bill_billing_entity"`
	BillBillingPeriodEndDate                                 *time.Time              `json:"bill_billing_period_end_date,omitempty" parquet:"name=bill_billing_period_end_date"`
	BillBillingPeriodStartDate                               *time.Time              `json:"bill_billing_period_start_date,omitempty" parquet:"name=bill_billing_period_start_date"`
	BillBillType                                             *string                 `json:"bill_bill_type,omitempty" parquet:"name=bill_bill_type"`
	BillInvoiceId                                            *string                 `json:"bill_invoice_id,omitempty" parquet:"name=bill_invoice_id"`
	BillInvoicingEntity                                      *string                 `json:"bill_invoicing_entity,omitempty" parquet:"name=bill_invoicing_entity"`
	BillPayerAccountId                                       *string                 `json:"bill_payer_account_id,omitempty" parquet:"name=bill_payer_account_id"`
	BillPayerAccountName                                     *string                 `json:"bill_payer_account_name,omitempty" parquet:"name=bill_payer_account_name"`
	CostCategory                                             *map[string]interface{} `json:"cost_category,omitempty" parquet:"name=cost_category"`
	CostCategoryProject                                      *string                 `json:"cost_category_project,omitempty" parquet:"name=cost_category_project"`
	CostCategoryTeam                                         *string                 `json:"cost_category_team,omitempty" parquet:"name=cost_category_team"`
	CostCategoryTest_Cost_Category                           *string                 `json:"cost_category_test_cost_category,omitempty" parquet:"name=cost_category_test_cost_category"`
	Discount                                                 *map[string]interface{} `json:"discount,omitempty" parquet:"name=discount"`
	DiscountBundledDiscount                                  *float64                `json:"discount_bundled_discount,omitempty" parquet:"name=discount_bundled_discount"`
	DiscountTotalDiscount                                    *float64                `json:"discount_total_discount,omitempty" parquet:"name=discount_total_discount"`
	IdentityLineItemId                                       *string                 `json:"identity_line_item_id,omitempty" parquet:"name=identity_line_item_id"`
	IdentityTimeInterval                                     *string                 `json:"identity_time_interval,omitempty" parquet:"name=identity_time_interval"`
	LineItemAvailabilityZone                                 *string                 `json:"line_item_availability_zone,omitempty" parquet:"name=line_item_availability_zone"`
	LineItemBlendedCost                                      *float64                `json:"line_item_blended_cost,omitempty" parquet:"name=line_item_blended_cost"`
	LineItemBlendedRate                                      *string                 `json:"line_item_blended_rate,omitempty" parquet:"name=line_item_blended_rate"`
	LineItemCurrencyCode                                     *string                 `json:"line_item_currency_code,omitempty" parquet:"name=line_item_currency_code"`
	LineItemLegalEntity                                      *string                 `json:"line_item_legal_entity,omitempty" parquet:"name=line_item_legal_entity"`
	LineItemLineItemDescription                              *string                 `json:"line_item_line_item_description,omitempty" parquet:"name=line_item_line_item_description"`
	LineItemLineItemType                                     *string                 `json:"line_item_line_item_type,omitempty" parquet:"name=line_item_line_item_type"`
	LineItemNetUnblendedCost                                 *float64                `json:"line_item_net_unblended_cost,omitempty" parquet:"name=line_item_net_unblended_cost"`
	LineItemNetUnblendedRate                                 *string                 `json:"line_item_net_unblended_rate,omitempty" parquet:"name=line_item_net_unblended_rate"`
	LineItemNormalizationFactor                              *float64                `json:"line_item_normalization_factor,omitempty" parquet:"name=line_item_normalization_factor"`
	LineItemNormalizedUsageAmount                            *float64                `json:"line_item_normalized_usage_amount,omitempty" parquet:"name=line_item_normalized_usage_amount"`
	LineItemOperation                                        *string                 `json:"line_item_operation,omitempty" parquet:"name=line_item_operation"`
	LineItemProductCode                                      *string                 `json:"line_item_product_code,omitempty" parquet:"name=line_item_product_code"`
	LineItemResourceId                                       *string                 `json:"line_item_resource_id,omitempty" parquet:"name=line_item_resource_id"`
	LineItemTaxType                                          *string                 `json:"line_item_tax_type,omitempty" parquet:"name=line_item_tax_type"`
	LineItemUnblendedCost                                    *float64                `json:"line_item_unblended_cost,omitempty" parquet:"name=line_item_unblended_cost"`
	LineItemUnblendedRate                                    *string                 `json:"line_item_unblended_rate,omitempty" parquet:"name=line_item_unblended_rate"`
	LineItemUsageAccountId                                   *string                 `json:"line_item_usage_account_id,omitempty" parquet:"name=line_item_usage_account_id"`
	LineItemUsageAccountName                                 *string                 `json:"line_item_usage_account_name,omitempty" parquet:"name=line_item_usage_account_name"`
	LineItemUsageAmount                                      *string                 `json:"line_item_usage_amount,omitempty" parquet:"name=line_item_usage_amount"`
	LineItemUsageEndDate                                     *time.Time              `json:"line_item_usage_end_date,omitempty" parquet:"name=line_item_usage_end_date"`
	LineItemUsageStartDate                                   *time.Time              `json:"line_item_usage_start_date,omitempty" parquet:"name=line_item_usage_start_date"`
	LineItemUsageType                                        *string                 `json:"line_item_usage_type,omitempty" parquet:"name=line_item_usage_type"`
	PricingCurrency                                          *string                 `json:"pricing_currency,omitempty" parquet:"name=pricing_currency"`
	PricingLeaseContractLength                               *string                 `json:"pricing_lease_contract_length,omitempty" parquet:"name=pricing_lease_contract_length"`
	PricingOfferingClass                                     *string                 `json:"pricing_offering_class,omitempty" parquet:"name=pricing_offering_class"`
	PricingPublicOnDemandCost                                *float64                `json:"pricing_public_on_demand_cost,omitempty" parquet:"name=pricing_public_on_demand_cost"`
	PricingPublicOnDemandRate                                *string                 `json:"pricing_public_on_demand_rate,omitempty" parquet:"name=pricing_public_on_demand_rate"`
	PricingPurchaseOption                                    *string                 `json:"pricing_purchase_option,omitempty" parquet:"name=pricing_purchase_option"`
	PricingRateCode                                          *string                 `json:"pricing_rate_code,omitempty" parquet:"name=pricing_rate_code"`
	PricingRateId                                            *string                 `json:"pricing_rate_id,omitempty" parquet:"name=pricing_rate_id"`
	PricingTerm                                              *string                 `json:"pricing_term,omitempty" parquet:"name=pricing_term"`
	PricingUnit                                              *string                 `json:"pricing_unit,omitempty" parquet:"name=pricing_unit"`
	Product                                                  *map[string]interface{} `json:"product,omitempty" parquet:"name=product"`
	ProductAlarmType                                         *string                 `json:"product_alarm_type,omitempty" parquet:"name=product_alarm_type"`
	ProductArchitecture                                      *string                 `json:"product_architecture,omitempty" parquet:"name=product_architecture"`
	ProductAttachmentType                                    *string                 `json:"product_attachment_type,omitempty" parquet:"name=product_attachment_type"`
	ProductAvailability                                      *string                 `json:"product_availability,omitempty" parquet:"name=product_availability"`
	ProductAvailabilityZone                                  *string                 `json:"product_availability_zone,omitempty" parquet:"name=product_availability_zone"`
	ProductBackupservice                                     *string                 `json:"product_backupservice,omitempty" parquet:"name=product_backupservice"`
	ProductBrokerEngine                                      *string                 `json:"product_broker_engine,omitempty" parquet:"name=product_broker_engine"`
	ProductCacheEngine                                       *string                 `json:"product_cache_engine,omitempty" parquet:"name=product_cache_engine"`
	ProductCacheType                                         *string                 `json:"product_cache_type,omitempty" parquet:"name=product_cache_type"`
	ProductCapacitystatus                                    *string                 `json:"product_capacitystatus,omitempty" parquet:"name=product_capacitystatus"`
	ProductCategory                                          *string                 `json:"product_category,omitempty" parquet:"name=product_category"`
	ProductCiType                                            *string                 `json:"product_ci_type,omitempty" parquet:"name=product_ci_type"`
	ProductClassicnetworkingsupport                          *string                 `json:"product_classicnetworkingsupport,omitempty" parquet:"name=product_classicnetworkingsupport"`
	ProductClockSpeed                                        *string                 `json:"product_clock_speed,omitempty" parquet:"name=product_clock_speed"`
	ProductCloudformationresourceProvider                    *string                 `json:"product_cloudformationresource_provider,omitempty" parquet:"name=product_cloudformationresource_provider"`
	ProductComment                                           *string                 `json:"product_comment,omitempty" parquet:"name=product_comment"`
	ProductComponent                                         *string                 `json:"product_component,omitempty" parquet:"name=product_component"`
	ProductComputeFamily                                     *string                 `json:"product_compute_family,omitempty" parquet:"name=product_compute_family"`
	ProductComputeType                                       *string                 `json:"product_compute_type,omitempty" parquet:"name=product_compute_type"`
	ProductContentType                                       *string                 `json:"product_content_type,omitempty" parquet:"name=product_content_type"`
	ProductCountsAgainstQuota                                *string                 `json:"product_counts_against_quota,omitempty" parquet:"name=product_counts_against_quota"`
	ProductCputype                                           *string                 `json:"product_cputype,omitempty" parquet:"name=product_cputype"`
	ProductCurrentGeneration                                 *string                 `json:"product_current_generation,omitempty" parquet:"name=product_current_generation"`
	ProductDatabaseEngine                                    *string                 `json:"product_database_engine,omitempty" parquet:"name=product_database_engine"`
	ProductDatabaseEngineType                                *string                 `json:"product_database_engine_type,omitempty" parquet:"name=product_database_engine_type"`
	ProductDatatransferout                                   *string                 `json:"product_datatransferout,omitempty" parquet:"name=product_datatransferout"`
	ProductDataTransferQuota                                 *string                 `json:"product_data_transfer_quota,omitempty" parquet:"name=product_data_transfer_quota"`
	ProductDedicatedEbsThroughput                            *string                 `json:"product_dedicated_ebs_throughput,omitempty" parquet:"name=product_dedicated_ebs_throughput"`
	ProductDeploymentOption                                  *string                 `json:"product_deployment_option,omitempty" parquet:"name=product_deployment_option"`
	ProductDescription                                       *string                 `json:"product_description,omitempty" parquet:"name=product_description"`
	ProductDirectorySize                                     *string                 `json:"product_directory_size,omitempty" parquet:"name=product_directory_size"`
	ProductDirectoryType                                     *string                 `json:"product_directory_type,omitempty" parquet:"name=product_directory_type"`
	ProductDirectoryTypeDescription                          *string                 `json:"product_directory_type_description,omitempty" parquet:"name=product_directory_type_description"`
	ProductDurability                                        *string                 `json:"product_durability,omitempty" parquet:"name=product_durability"`
	ProductEcu                                               *string                 `json:"product_ecu,omitempty" parquet:"name=product_ecu"`
	ProductEdition                                           *string                 `json:"product_edition,omitempty" parquet:"name=product_edition"`
	ProductEndpoint                                          *string                 `json:"product_endpoint,omitempty" parquet:"name=product_endpoint"`
	ProductEndpointType                                      *string                 `json:"product_endpoint_type,omitempty" parquet:"name=product_endpoint_type"`
	ProductEngine                                            *string                 `json:"product_engine,omitempty" parquet:"name=product_engine"`
	ProductEngineCode                                        *string                 `json:"product_engine_code,omitempty" parquet:"name=product_engine_code"`
	ProductEngineMajorVersion                                *string                 `json:"product_engine_major_version,omitempty" parquet:"name=product_engine_major_version"`
	ProductEnhancedNetworkingSupport                         *string                 `json:"product_enhanced_networking_support,omitempty" parquet:"name=product_enhanced_networking_support"`
	ProductEnhancedNetworkingSupported                       *string                 `json:"product_enhanced_networking_supported,omitempty" parquet:"name=product_enhanced_networking_supported"`
	ProductEventType                                         *string                 `json:"product_event_type,omitempty" parquet:"name=product_event_type"`
	ProductExtendedSupportPricingYear                        *string                 `json:"product_extended_support_pricing_year,omitempty" parquet:"name=product_extended_support_pricing_year"`
	ProductFeeCode                                           *string                 `json:"product_fee_code,omitempty" parquet:"name=product_fee_code"`
	ProductFeeDescription                                    *string                 `json:"product_fee_description,omitempty" parquet:"name=product_fee_description"`
	ProductFileSystemType                                    *string                 `json:"product_file_system_type,omitempty" parquet:"name=product_file_system_type"`
	ProductFindingGroup                                      *string                 `json:"product_finding_group,omitempty" parquet:"name=product_finding_group"`
	ProductFindingSource                                     *string                 `json:"product_finding_source,omitempty" parquet:"name=product_finding_source"`
	ProductFindingStorage                                    *string                 `json:"product_finding_storage,omitempty" parquet:"name=product_finding_storage"`
	ProductFreeOverage                                       *string                 `json:"product_free_overage,omitempty" parquet:"name=product_free_overage"`
	ProductFreeUsageIncluded                                 *string                 `json:"product_free_usage_included,omitempty" parquet:"name=product_free_usage_included"`
	ProductFromLocation                                      *string                 `json:"product_from_location,omitempty" parquet:"name=product_from_location"`
	ProductFromLocationType                                  *string                 `json:"product_from_location_type,omitempty" parquet:"name=product_from_location_type"`
	ProductFromRegionCode                                    *string                 `json:"product_from_region_code,omitempty" parquet:"name=product_from_region_code"`
	ProductGb                                                *string                 `json:"product_gb,omitempty" parquet:"name=product_gb"`
	ProductGpu                                               *string                 `json:"product_gpu,omitempty" parquet:"name=product_gpu"`
	ProductGpuMemory                                         *string                 `json:"product_gpu_memory,omitempty" parquet:"name=product_gpu_memory"`
	ProductGraphqloperation                                  *string                 `json:"product_graphqloperation,omitempty" parquet:"name=product_graphqloperation"`
	ProductGroup                                             *string                 `json:"product_group,omitempty" parquet:"name=product_group"`
	ProductGroupDescription                                  *string                 `json:"product_group_description,omitempty" parquet:"name=product_group_description"`
	ProductIndexingSource                                    *string                 `json:"product_indexing_source,omitempty" parquet:"name=product_indexing_source"`
	ProductInsightstype                                      *string                 `json:"product_insightstype,omitempty" parquet:"name=product_insightstype"`
	ProductInstance                                          *string                 `json:"product_instance,omitempty" parquet:"name=product_instance"`
	ProductInstanceConfigurationType                         *string                 `json:"product_instance_configuration_type,omitempty" parquet:"name=product_instance_configuration_type"`
	ProductInstanceFamily                                    *string                 `json:"product_instance_family,omitempty" parquet:"name=product_instance_family"`
	ProductInstanceName                                      *string                 `json:"product_instance_name,omitempty" parquet:"name=product_instance_name"`
	ProductInstancesku                                       *string                 `json:"product_instancesku,omitempty" parquet:"name=product_instancesku"`
	ProductInstanceType                                      *string                 `json:"product_instance_type,omitempty" parquet:"name=product_instance_type"`
	ProductInstanceTypeFamily                                *string                 `json:"product_instance_type_family,omitempty" parquet:"name=product_instance_type_family"`
	ProductIntelAvx2Available                                *string                 `json:"product_intel_avx2_available,omitempty" parquet:"name=product_intel_avx2_available"`
	ProductIntelAvx2available                                *string                 `json:"product_intel_avx2available,omitempty" parquet:"name=product_intel_avx2available"`
	ProductIntelAvxAvailable                                 *string                 `json:"product_intel_avx_available,omitempty" parquet:"name=product_intel_avx_available"`
	ProductIntelTurboAvailable                               *string                 `json:"product_intel_turbo_available,omitempty" parquet:"name=product_intel_turbo_available"`
	ProductIo                                                *string                 `json:"product_io,omitempty" parquet:"name=product_io"`
	ProductLicenseModel                                      *string                 `json:"product_license_model,omitempty" parquet:"name=product_license_model"`
	ProductLocation                                          *string                 `json:"product_location,omitempty" parquet:"name=product_location"`
	ProductLocationType                                      *string                 `json:"product_location_type,omitempty" parquet:"name=product_location_type"`
	ProductLogsDestination                                   *string                 `json:"product_logs_destination,omitempty" parquet:"name=product_logs_destination"`
	ProductMarketoption                                      *string                 `json:"product_marketoption,omitempty" parquet:"name=product_marketoption"`
	ProductMaximumExtendedStorage                            *string                 `json:"product_maximum_extended_storage,omitempty" parquet:"name=product_maximum_extended_storage"`
	ProductMaxIopsBurstPerformance                           *string                 `json:"product_max_iops_burst_performance,omitempty" parquet:"name=product_max_iops_burst_performance"`
	ProductMaxIopsvolume                                     *string                 `json:"product_max_iopsvolume,omitempty" parquet:"name=product_max_iopsvolume"`
	ProductMaxThroughputvolume                               *string                 `json:"product_max_throughputvolume,omitempty" parquet:"name=product_max_throughputvolume"`
	ProductMaxVolumeSize                                     *string                 `json:"product_max_volume_size,omitempty" parquet:"name=product_max_volume_size"`
	ProductMemory                                            *string                 `json:"product_memory,omitempty" parquet:"name=product_memory"`
	ProductMemoryGib                                         *string                 `json:"product_memory_gib,omitempty" parquet:"name=product_memory_gib"`
	ProductMemorytype                                        *string                 `json:"product_memorytype,omitempty" parquet:"name=product_memorytype"`
	ProductMessageDeliveryFrequency                          *string                 `json:"product_message_delivery_frequency,omitempty" parquet:"name=product_message_delivery_frequency"`
	ProductMessageDeliveryOrder                              *string                 `json:"product_message_delivery_order,omitempty" parquet:"name=product_message_delivery_order"`
	ProductMeteringType                                      *string                 `json:"product_metering_type,omitempty" parquet:"name=product_metering_type"`
	ProductMinVolumeSize                                     *string                 `json:"product_min_volume_size,omitempty" parquet:"name=product_min_volume_size"`
	ProductNetworkPerformance                                *string                 `json:"product_network_performance,omitempty" parquet:"name=product_network_performance"`
	ProductNewcode                                           *string                 `json:"product_newcode,omitempty" parquet:"name=product_newcode"`
	ProductNormalizationSizeFactor                           *string                 `json:"product_normalization_size_factor,omitempty" parquet:"name=product_normalization_size_factor"`
	ProductOperatingSystem                                   *string                 `json:"product_operating_system,omitempty" parquet:"name=product_operating_system"`
	ProductOperation                                         *string                 `json:"product_operation,omitempty" parquet:"name=product_operation"`
	ProductOrigin                                            *string                 `json:"product_origin,omitempty" parquet:"name=product_origin"`
	ProductOverageType                                       *string                 `json:"product_overage_type,omitempty" parquet:"name=product_overage_type"`
	ProductOverhead                                          *string                 `json:"product_overhead,omitempty" parquet:"name=product_overhead"`
	ProductPackSize                                          *string                 `json:"product_pack_size,omitempty" parquet:"name=product_pack_size"`
	ProductParameterType                                     *string                 `json:"product_parameter_type,omitempty" parquet:"name=product_parameter_type"`
	ProductPhysicalCpu                                       *string                 `json:"product_physical_cpu,omitempty" parquet:"name=product_physical_cpu"`
	ProductPhysicalGpu                                       *string                 `json:"product_physical_gpu,omitempty" parquet:"name=product_physical_gpu"`
	ProductPhysicalProcessor                                 *string                 `json:"product_physical_processor,omitempty" parquet:"name=product_physical_processor"`
	ProductPlatoclassificationtype                           *string                 `json:"product_platoclassificationtype,omitempty" parquet:"name=product_platoclassificationtype"`
	ProductPlatoinstancename                                 *string                 `json:"product_platoinstancename,omitempty" parquet:"name=product_platoinstancename"`
	ProductPlatoinstancetype                                 *string                 `json:"product_platoinstancetype,omitempty" parquet:"name=product_platoinstancetype"`
	ProductPlatopricingtype                                  *string                 `json:"product_platopricingtype,omitempty" parquet:"name=product_platopricingtype"`
	ProductPreInstalledSw                                    *string                 `json:"product_pre_installed_sw,omitempty" parquet:"name=product_pre_installed_sw"`
	ProductPricingplan                                       *string                 `json:"product_pricingplan,omitempty" parquet:"name=product_pricingplan"`
	ProductPricingUnit                                       *string                 `json:"product_pricing_unit,omitempty" parquet:"name=product_pricing_unit"`
	ProductProcessorArchitecture                             *string                 `json:"product_processor_architecture,omitempty" parquet:"name=product_processor_architecture"`
	ProductProcessorFeatures                                 *string                 `json:"product_processor_features,omitempty" parquet:"name=product_processor_features"`
	ProductProductFamily                                     *string                 `json:"product_product_family,omitempty" parquet:"name=product_product_family"`
	ProductProductName                                       *string                 `json:"product_product_name,omitempty" parquet:"name=product_product_name"`
	ProductProtocol                                          *string                 `json:"product_protocol,omitempty" parquet:"name=product_protocol"`
	ProductProvider                                          *string                 `json:"product_provider,omitempty" parquet:"name=product_provider"`
	ProductProvisioned                                       *string                 `json:"product_provisioned,omitempty" parquet:"name=product_provisioned"`
	ProductQPresent                                          *string                 `json:"product_q_present,omitempty" parquet:"name=product_q_present"`
	ProductQueueType                                         *string                 `json:"product_queue_type,omitempty" parquet:"name=product_queue_type"`
	ProductRecipient                                         *string                 `json:"product_recipient,omitempty" parquet:"name=product_recipient"`
	ProductRegion                                            *string                 `json:"product_region,omitempty" parquet:"name=product_region"`
	ProductRegionCode                                        *string                 `json:"product_region_code,omitempty" parquet:"name=product_region_code"`
	ProductRequestDescription                                *string                 `json:"product_request_description,omitempty" parquet:"name=product_request_description"`
	ProductRequestType                                       *string                 `json:"product_request_type,omitempty" parquet:"name=product_request_type"`
	ProductResourceEndpoint                                  *string                 `json:"product_resource_endpoint,omitempty" parquet:"name=product_resource_endpoint"`
	ProductResourcePriceGroup                                *string                 `json:"product_resource_price_group,omitempty" parquet:"name=product_resource_price_group"`
	ProductRoutingTarget                                     *string                 `json:"product_routing_target,omitempty" parquet:"name=product_routing_target"`
	ProductRoutingType                                       *string                 `json:"product_routing_type,omitempty" parquet:"name=product_routing_type"`
	ProductScanType                                          *string                 `json:"product_scan_type,omitempty" parquet:"name=product_scan_type"`
	ProductServicecode                                       *string                 `json:"product_servicecode,omitempty" parquet:"name=product_servicecode"`
	ProductServicename                                       *string                 `json:"product_servicename,omitempty" parquet:"name=product_servicename"`
	ProductSizeFlex                                          *string                 `json:"product_size_flex,omitempty" parquet:"name=product_size_flex"`
	ProductSku                                               *string                 `json:"product_sku,omitempty" parquet:"name=product_sku"`
	ProductSoftwareType                                      *string                 `json:"product_software_type,omitempty" parquet:"name=product_software_type"`
	ProductStandardGroup                                     *string                 `json:"product_standard_group,omitempty" parquet:"name=product_standard_group"`
	ProductStandardStorage                                   *string                 `json:"product_standard_storage,omitempty" parquet:"name=product_standard_storage"`
	ProductStandardStorageRetentionIncluded                  *string                 `json:"product_standard_storage_retention_included,omitempty" parquet:"name=product_standard_storage_retention_included"`
	ProductSteps                                             *string                 `json:"product_steps,omitempty" parquet:"name=product_steps"`
	ProductStorage                                           *string                 `json:"product_storage,omitempty" parquet:"name=product_storage"`
	ProductStorageCapacityQuota                              *string                 `json:"product_storage_capacity_quota,omitempty" parquet:"name=product_storage_capacity_quota"`
	ProductStorageClass                                      *string                 `json:"product_storage_class,omitempty" parquet:"name=product_storage_class"`
	ProductStorageFamily                                     *string                 `json:"product_storage_family,omitempty" parquet:"name=product_storage_family"`
	ProductStorageMedia                                      *string                 `json:"product_storage_media,omitempty" parquet:"name=product_storage_media"`
	ProductStorageTier                                       *string                 `json:"product_storage_tier,omitempty" parquet:"name=product_storage_tier"`
	ProductStorageType                                       *string                 `json:"product_storage_type,omitempty" parquet:"name=product_storage_type"`
	ProductSubcategory                                       *string                 `json:"product_subcategory,omitempty" parquet:"name=product_subcategory"`
	ProductSubscriptionType                                  *string                 `json:"product_subscription_type,omitempty" parquet:"name=product_subscription_type"`
	ProductSubservice                                        *string                 `json:"product_subservice,omitempty" parquet:"name=product_subservice"`
	ProductTenancy                                           *string                 `json:"product_tenancy,omitempty" parquet:"name=product_tenancy"`
	ProductThroughput                                        *string                 `json:"product_throughput,omitempty" parquet:"name=product_throughput"`
	ProductThroughputCapacity                                *string                 `json:"product_throughput_capacity,omitempty" parquet:"name=product_throughput_capacity"`
	ProductTiertype                                          *string                 `json:"product_tiertype,omitempty" parquet:"name=product_tiertype"`
	ProductToLocation                                        *string                 `json:"product_to_location,omitempty" parquet:"name=product_to_location"`
	ProductToLocationType                                    *string                 `json:"product_to_location_type,omitempty" parquet:"name=product_to_location_type"`
	ProductToRegionCode                                      *string                 `json:"product_to_region_code,omitempty" parquet:"name=product_to_region_code"`
	ProductTransferType                                      *string                 `json:"product_transfer_type,omitempty" parquet:"name=product_transfer_type"`
	ProductType                                              *string                 `json:"product_type,omitempty" parquet:"name=product_type"`
	ProductUsageFamily                                       *string                 `json:"product_usage_family,omitempty" parquet:"name=product_usage_family"`
	ProductUsageGroup                                        *string                 `json:"product_usage_group,omitempty" parquet:"name=product_usage_group"`
	ProductUsagetype                                         *string                 `json:"product_usagetype,omitempty" parquet:"name=product_usagetype"`
	ProductUsageVolume                                       *string                 `json:"product_usage_volume,omitempty" parquet:"name=product_usage_volume"`
	ProductVaulttype                                         *string                 `json:"product_vaulttype,omitempty" parquet:"name=product_vaulttype"`
	ProductVcpu                                              *string                 `json:"product_vcpu,omitempty" parquet:"name=product_vcpu"`
	ProductVersion                                           *string                 `json:"product_version,omitempty" parquet:"name=product_version"`
	ProductVolumeApiName                                     *string                 `json:"product_volume_api_name,omitempty" parquet:"name=product_volume_api_name"`
	ProductVolumeType                                        *string                 `json:"product_volume_type,omitempty" parquet:"name=product_volume_type"`
	ProductVpcnetworkingsupport                              *string                 `json:"product_vpcnetworkingsupport,omitempty" parquet:"name=product_vpcnetworkingsupport"`
	ProductWithActiveUsers                                   *string                 `json:"product_with_active_users,omitempty" parquet:"name=product_with_active_users"`
	Reservation                                              *map[string]interface{} `json:"reservation,omitempty" parquet:"name=reservation"`
	ReservationAmortizedUpfrontCostForUsage                  *float64                `json:"reservation_amortized_upfront_cost_for_usage,omitempty" parquet:"name=reservation_amortized_upfront_cost_for_usage"`
	ReservationAmortizedUpfrontFeeForBillingPeriod           *float64                `json:"reservation_amortized_upfront_fee_for_billing_period,omitempty" parquet:"name=reservation_amortized_upfront_fee_for_billing_period"`
	ReservationArn                                           *string                 `json:"reservation_reservation_arn,omitempty" parquet:"name=reservation_reservation_arn"`
	ReservationAvailabilityZone                              *string                 `json:"reservation_availability_zone,omitempty" parquet:"name=reservation_availability_zone"`
	ReservationEffectiveCost                                 *float64                `json:"reservation_effective_cost,omitempty" parquet:"name=reservation_effective_cost"`
	ReservationEndTime                                       *string                 `json:"reservation_end_time,omitempty" parquet:"name=reservation_end_time"`
	ReservationModificationStatus                            *string                 `json:"reservation_modification_status,omitempty" parquet:"name=reservation_modification_status"`
	ReservationNetAmortizedUpfrontCostForUsage               *float64                `json:"reservation_net_amortized_upfront_cost_for_usage,omitempty" parquet:"name=reservation_net_amortized_upfront_cost_for_usage"`
	ReservationNetAmortizedUpfrontFeeForBillingPeriod        *float64                `json:"reservation_net_amortized_upfront_fee_for_billing_period,omitempty" parquet:"name=reservation_net_amortized_upfront_fee_for_billing_period"`
	ReservationNetEffectiveCost                              *float64                `json:"reservation_net_effective_cost,omitempty" parquet:"name=reservation_net_effective_cost"`
	ReservationNetRecurringFeeForUsage                       *float64                `json:"reservation_net_recurring_fee_for_usage,omitempty" parquet:"name=reservation_net_recurring_fee_for_usage"`
	ReservationNetUnusedAmortizedUpfrontFeeForBillingPeriod  *float64                `json:"reservation_net_unused_amortized_upfront_fee_for_billing_period,omitempty" parquet:"name=reservation_net_unused_amortized_upfront_fee_for_billing_period"`
	ReservationNetUnusedRecurringFee                         *float64                `json:"reservation_net_unused_recurring_fee,omitempty" parquet:"name=reservation_net_unused_recurring_fee"`
	ReservationNetUpfrontValue                               *float64                `json:"reservation_net_upfront_value,omitempty" parquet:"name=reservation_net_upfront_value"`
	ReservationNormalizedUnitsPerReservation                 *string                 `json:"reservation_normalized_units_per_reservation,omitempty" parquet:"name=reservation_normalized_units_per_reservation"`
	ReservationNumberOfReservations                          *string                 `json:"reservation_number_of_reservations,omitempty" parquet:"name=reservation_number_of_reservations"`
	ReservationRecurringFeeForUsage                          *float64                `json:"reservation_recurring_fee_for_usage,omitempty" parquet:"name=reservation_recurring_fee_for_usage"`
	ReservationStartTime                                     *time.Time              `json:"reservation_start_time,omitempty" parquet:"name=reservation_start_time"`
	ReservationSubscriptionId                                *string                 `json:"reservation_subscription_id,omitempty" parquet:"name=reservation_subscription_id"`
	ReservationTotalReservedNormalizedUnits                  *string                 `json:"reservation_total_reserved_normalized_units,omitempty" parquet:"name=reservation_total_reserved_normalized_units"`
	ReservationTotalReservedUnits                            *string                 `json:"reservation_total_reserved_units,omitempty" parquet:"name=reservation_total_reserved_units"`
	ReservationUnitsPerReservation                           *string                 `json:"reservation_units_per_reservation,omitempty" parquet:"name=reservation_units_per_reservation"`
	ReservationUnusedAmortizedUpfrontFeeForBillingPeriod     *float64                `json:"reservation_unused_amortized_upfront_fee_for_billing_period,omitempty" parquet:"name=reservation_unused_amortized_upfront_fee_for_billing_period"`
	ReservationUnusedNormalizedUnitQuantity                  *float64                `json:"reservation_unused_normalized_unit_quantity,omitempty" parquet:"name=reservation_unused_normalized_unit_quantity"`
	ReservationUnusedQuantity                                *float64                `json:"reservation_unused_quantity,omitempty" parquet:"name=reservation_unused_quantity"`
	ReservationUnusedRecurringFee                            *float64                `json:"reservation_unused_recurring_fee,omitempty" parquet:"name=reservation_unused_recurring_fee"`
	ReservationUpfrontValue                                  *int64                  `json:"reservation_upfront_value,omitempty" parquet:"name=reservation_upfront_value"`
	ResourceTags                                             *map[string]interface{} `json:"resource_tags,omitempty" parquet:"name=resource_tags"`
	ResourceTagsAwsCloudformationStackName                   *string                 `json:"resource_tags_aws_cloudformation_stack-name,omitempty" parquet:"name=resource_tags_aws_cloudformation_stack_name"`
	ResourceTagsAwsCreatedBy                                 *string                 `json:"resource_tags_aws_created_by,omitempty" parquet:"name=resource_tags_aws_created_by"`
	ResourceTagsUserCostCenter                               *string                 `json:"resource_tags_user_cost Center,omitempty" parquet:"name=resource_tags_user_cost_center"`
	ResourceTagsUserCreatedByActor                           *string                 `json:"resource_tags_user_created_by_actor,omitempty" parquet:"name=resource_tags_user_created_by_actor"`
	ResourceTagsUserDepartment                               *string                 `json:"resource_tags_user_department,omitempty" parquet:"name=resource_tags_user_department"`
	ResourceTagsUserName                                     *string                 `json:"resource_tags_user_name,omitempty" parquet:"name=resource_tags_user_name"`
	ResourceTagsUserOwner                                    *string                 `json:"resource_tags_user_owner,omitempty" parquet:"name=resource_tags_user_owner"`
	SavingsPlanAmortizedUpfrontCommitmentForBillingPeriod    *float64                `json:"savings_plan_amortized_upfront_commitment_for_billing_period,omitempty" parquet:"name=savings_plan_amortized_upfront_commitment_for_billing_period"`
	SavingsPlanEndTime                                       *time.Time              `json:"savings_plan_end_time,omitempty" parquet:"name=savings_plan_end_time"`
	SavingsPlanInstanceTypeFamily                            *string                 `json:"savings_plan_instance_type_family,omitempty" parquet:"name=savings_plan_instance_type_family"`
	SavingsPlanNetAmortizedUpfrontCommitmentForBillingPeriod *float64                `json:"savings_plan_net_amortized_upfront_commitment_for_billing_period,omitempty" parquet:"name=savings_plan_net_amortized_upfront_commitment_for_billing_period"`
	SavingsPlanNetRecurringCommitmentForBillingPeriod        *float64                `json:"savings_plan_net_recurring_commitment_for_billing_period,omitempty" parquet:"name=savings_plan_net_recurring_commitment_for_billing_period"`
	SavingsPlanNetSavingsPlanEffectiveCost                   *float64                `json:"savings_plan_net_savings_plan_effective_cost,omitempty" parquet:"name=savings_plan_net_savings_plan_effective_cost"`
	SavingsPlanOfferingType                                  *string                 `json:"savings_plan_offering_type,omitempty" parquet:"name=savings_plan_offering_type"`
	SavingsPlanPaymentOption                                 *string                 `json:"savings_plan_payment_option,omitempty" parquet:"name=savings_plan_payment_option"`
	SavingsPlanPurchaseTerm                                  *string                 `json:"savings_plan_purchase_term,omitempty" parquet:"name=savings_plan_purchase_term"`
	SavingsPlanRecurringCommitmentForBillingPeriod           *float64                `json:"savings_plan_recurring_commitment_for_billing_period,omitempty" parquet:"name=savings_plan_recurring_commitment_for_billing_period"`
	SavingsPlanRegion                                        *string                 `json:"savings_plan_region,omitempty" parquet:"name=savings_plan_region"`
	SavingsPlanSavingsPlanARN                                *string                 `json:"savings_plan_savings_plan_arn,omitempty" parquet:"name=savings_plan_savings_plan_arn"`
	SavingsPlanSavingsPlanEffectiveCost                      *float64                `json:"savings_plan_savings_plan_effective_cost,omitempty" parquet:"name=savings_plan_savings_plan_effective_cost"`
	SavingsPlanSavingsPlanRate                               *string                 `json:"savings_plan_savings_plan_rate,omitempty" parquet:"name=savings_plan_savings_plan_rate"`
	SavingsPlanStartTime                                     *time.Time              `json:"savings_plan_start_time,omitempty" parquet:"name=savings_plan_start_time"`
	SavingsPlanTotalCommitmentToDate                         *string                 `json:"savings_plan_total_commitment_to_date,omitempty" parquet:"name=savings_plan_total_commitment_to_date"`
	SavingsPlanUsedCommitment                                *string                 `json:"savings_plan_used_commitment,omitempty" parquet:"name=savings_plan_used_commitment"`
	SplitLineItemActualUsage                                 *float64                `json:"split_line_item_actual_usage,omitempty" parquet:"name=split_line_item_actual_usage"`
	SplitLineItemNetSplitCost                                *float64                `json:"split_line_item_net_split_cost,omitempty" parquet:"name=split_line_item_net_split_cost"`
	SplitLineItemNetUnusedCost                               *float64                `json:"split_line_item_net_unused_cost,omitempty" parquet:"name=split_line_item_net_unused_cost"`
	SplitLineItemParentResourceId                            *string                 `json:"split_line_item_parent_resource_id,omitempty" parquet:"name=split_line_item_parent_resource_id"`
	SplitLineItemPublicOnDemandSplitCost                     *float64                `json:"split_line_item_public_on_demand_split_cost,omitempty" parquet:"name=split_line_item_public_on_demand_split_cost"`
	SplitLineItemPublicOnDemandUnusedCost                    *float64                `json:"split_line_item_public_on_demand_unused_cost,omitempty" parquet:"name=split_line_item_public_on_demand_unused_cost"`
	SplitLineItemReservedUsage                               *float64                `json:"split_line_item_reserved_usage,omitempty" parquet:"name=split_line_item_reserved_usage"`
	SplitLineItemSplitCost                                   *float64                `json:"split_line_item_split_cost,omitempty" parquet:"name=split_line_item_split_cost"`
	SplitLineItemSplitUsage                                  *float64                `json:"split_line_item_split_usage,omitempty" parquet:"name=split_line_item_split_usage"`
	SplitLineItemSplitUsageRatio                             *string                 `json:"split_line_item_split_usage_ratio,omitempty" parquet:"name=split_line_item_split_usage_ratio"`
	SplitLineItemUnusedCost                                  *float64                `json:"split_line_item_unused_cost,omitempty" parquet:"name=split_line_item_unused_cost"`
}

func (c *CostAndUsageLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"bill_billing_entity":                                      "Helps in identify whether the invoices or transactions are for AWS Marketplace or for purchases of other AWS services.",
		"bill_billing_period_end_date":                             "The end date of the billing period that is covered by this report, in UTC. The format is YYYY-MM-DDTHH:mm:ssZ.",
		"bill_billing_period_start_date":                           "The start date of the billing period that is covered by this report, in UTC. The format is YYYY-MM-DDTHH:mm:ssZ.",
		"bill_bill_type":                                           "The type of bill that this report covers.",
		"bill_invoice_id":                                          "The ID associated with a specific line item. Until the report is final, the InvoiceId is blank.",
		"bill_invoicing_entity":                                    "The AWS entity that issues the invoice.",
		"bill_payer_account_id":                                    "The account ID of the paying account. For an organization in AWS Organizations, this is the account ID of the management account.",
		"bill_payer_account_name":                                  "The account name of the paying account. For an organization in AWS Organizations, this is the name of the management account.",
		"cost_category":                                            "Cost Category entries are automatically populated when you create a Cost Category and categorization rule. These entries include user-defined Cost Category names as keys, and corresponding Cost Category values.",
		"cost_category_project":                                    "",
		"cost_category_team":                                         "",
		"cost_category_test_cost_category":                         "",
		"discount":                                                 "A map column that contains key-value pairs of additional discount data for a given line item when applicable.",
		"discount_bundled_discount":                                "The bundled discount applied to the line item. A bundled discount is a usage-based discount that provides free or discounted usage of a service or feature based on the usage of another service or feature.",
		"discount_total_discount":                                  "The sum of all the discount columns for the corresponding line item.",
		"identity_line_item_id":                                    "This field is generated for each line item and is unique in a given partition. This does not guarantee that the field will be unique across an entire delivery (that is, all partitions in an update) of the AWS CUR. The line item ID isn't consistent between different Cost and Usage Reports and can't be used to identify the same line item across different reports.",
		"identity_time_interval":                                   "The time interval that this line item applies to, in the following format: YYYY-MM-DDTHH:mm:ssZ/YYYY-MM-DDTHH:mm:ssZ. The time interval is in UTC and can be either daily or hourly, depending on the granularity of the report.",
		"line_item_availability_zone":                              "The Availability Zone that hosts this line item.",
		"line_item_blended_cost":                                   "The BlendedRate multiplied by the UsageAmount.",
		"line_item_blended_rate":                                   "The BlendedRate is the average cost incurred for each SKU across an organization.",
		"line_item_currency_code":                                  "The currency that this line item is shown in. All AWS customers are billed in US dollars by default. To change your billing currency, see Changing which currency you use to pay your bill in the AWS Billing User Guide.",
		"line_item_legal_entity":                                   "The Seller of Record of a specific product or service. In most cases, the invoicing entity and legal entity are the same. The values might differ for third-party AWS Marketplace transactions.",
		"line_item_line_item_description":                          "The description of the line item type.",
		"line_item_line_item_type":                                 "The type of charge covered by this line item.",
		"line_item_net_unblended_cost":                             "The actual after-discount cost that you're paying for the line item.",
		"line_item_net_unblended_rate":                             "The actual after-discount rate that you're paying for the line item.",
		"line_item_normalization_factor":                           "As long as the instance has shared tenancy, AWS can apply all Regional Linux or Unix Amazon EC2 and Amazon RDS RI discounts to all instance sizes in an instance family and AWS Region. This also applies to RI discounts for member accounts in an organization. All new and existing Amazon EC2 and Amazon RDS size-flexible RIs are sized according to a normalization factor, based on the instance size.",
		"line_item_normalized_usage_amount":                        "The amount of usage that you incurred, in normalized units, for size-flexible RIs. The NormalizedUsageAmount is equal to UsageAmount multiplied by NormalizationFactor.",
		"line_item_operation":                                      "The specific AWS operation covered by this line item. This describes the specific usage of the line item.",
		"line_item_product_code":                                   "The code of the product measured.",
		"line_item_resource_id":                                    "If you chose to include individual resource IDs in your report, this column contains the ID of the resource that you provisioned.",
		"line_item_tax_type":                                       "The type of tax that AWS applied to this line item.",
		"line_item_unblended_cost":                                 "The UnblendedCost is the UnblendedRate multiplied by the UsageAmount.",
		"line_item_unblended_rate":                                 "In consolidated billing for accounts using AWS Organizations, the unblended rate is the rate associated with an individual account's service usage. For Amazon EC2 and Amazon RDS line items that have an RI discount applied to them, the UnblendedRate is zero. Line items with an RI discount have a LineItemType of DiscountedUsage.",
		"line_item_usage_account_id":                               "The account ID of the account that used this line item. For organizations, this can be either the management account or a member account. You can use this field to track costs or usage by account.",
		"line_item_usage_account_name":                             "The name of the account that used this line item. For organizations, this can be either the management account or a member account. You can use this field to track costs or usage by account.",
		"line_item_usage_amount":                                   "The amount of usage that you incurred during the specified time period. For size-flexible Reserved Instances, use the reservation_total_reserved_units column instead. Certain subscription charges will have a UsageAmount of 0.",
		"line_item_usage_end_date":                                 "The end date and time for the line item in UTC, exclusive. The format is YYYY-MM-DDTHH:mm:ssZ.",
		"line_item_usage_start_date":                               "The start date and time for the line item in UTC, inclusive. The format is YYYY-MM-DDTHH:mm:ssZ.",
		"line_item_usage_type":                                     "The usage details of the line item.",
		"pricing_currency":                                         "The currency that the pricing data is shown in.",
		"pricing_lease_contract_length":                            "The length of time that your RI is reserved for.",
		"pricing_offering_class":                                   "The offering class of the Reserved Instance.",
		"pricing_public_on_demand_cost":                            "The total cost for the line item based on public On-Demand Instance rates. If you have SKUs with multiple On-Demand public costs, the equivalent cost for the highest tier is displayed. For example, services offering free-tiers or tiered pricing.",
		"pricing_public_on_demand_rate":                            "The public On-Demand Instance rate in this billing period for the specific line item of usage. If you have SKUs with multiple On-Demand public rates, the equivalent rate for the highest tier is displayed. For example, services offering free-tiers or tiered pricing.",
		"pricing_purchase_option":                                  "How you chose to pay for this line item. Valid values are All Upfront, Partial Upfront, and No Upfront.",
		"pricing_rate_code":                                        "A unique code for a product/ offer/ pricing-tier combination. The product and term combinations can have multiple price dimensions, such as a free tier, low-use tier, and high-use tier.",
		"pricing_rate_id":                                          "The ID of the rate for a line item.",
		"pricing_term":                                             "Whether your AWS usage is Reserved or On-Demand.",
		"pricing_unit":                                             "The pricing unit that AWS used for calculating your usage cost. For example, the pricing unit for Amazon EC2 instance usage is in hours.",
		"product":                                                  "A map column for where each key-value pair is an additional product attribute and its value.",
		"product_alarm_type":                                       "",
		"product_architecture":                                     "",
		"product_attachment_type":                                  "",
		"product_availability":                                     "",
		"product_availability_zone":                                "",
		"product_backupservice":                                    "",
		"product_broker_engine":                                    "",
		"product_cache_engine":                                     "",
		"product_cache_type":                                       "",
		"product_capacitystatus":                                   "",
		"product_category":                                         "",
		"product_ci_type":                                          "",
		"product_classicnetworkingsupport":                         "",
		"product_clock_speed":                                      "",
		"product_cloudformationresource_provider":                  "",
		"product_comment":                                          "A comment regarding the product.",
		"product_component":                                        "",
		"product_compute_family":                                   "",
		"product_compute_type":                                     "",
		"product_content_type":                                     "",
		"product_counts_against_quota":                             "",
		"product_cputype":                                          "",
		"product_current_generation":                               "",
		"product_database_engine":                                  "",
		"product_database_engine_type":                             "",
		"product_datatransferout":                                  "",
		"product_data_transfer_quota":                              "",
		"product_dedicated_ebs_throughput":                         "",
		"product_deployment_option":                                "",
		"product_description":                                      "",
		"product_directory_size":                                   "",
		"product_directory_type":                                   "",
		"product_directory_type_description":                       "",
		"product_durability":                                       "",
		"product_ecu":                                              "",
		"product_edition":                                          "",
		"product_endpoint":                                         "",
		"product_endpoint_type":                                    "",
		"product_engine":                                           "",
		"product_engine_code":                                      "",
		"product_engine_major_version":                             "",
		"product_enhanced_networking_support":                      "",
		"product_enhanced_networking_supported":                    "",
		"product_event_type":                                       "",
		"product_extended_support_pricing_year":                    "",
		"product_fee_code":                                         "",
		"product_fee_description":                                  "",
		"product_file_system_type":                                 "",
		"product_finding_group":                                    "",
		"product_finding_source":                                   "",
		"product_finding_storage":                                  "",
		"product_free_overage":                                     "",
		"product_free_usage_included":                              "",
		"product_from_location":                                    "",
		"product_from_location_type":                               "",
		"product_from_region_code":                                 "",
		"product_gb":                                               "",
		"product_gpu":                                              "",
		"product_gpu_memory":                                       "",
		"product_graphqloperation":                                 "",
		"product_group":                                            "",
		"product_group_description":                                "",
		"product_indexing_source":                                  "",
		"product_insightstype":                                     "",
		"product_instance":                                         "",
		"product_instance_configuration_type":                      "",
		"product_instance_family":                                  "",
		"product_instance_name":                                    "",
		"product_instancesku":                                      "",
		"product_instance_type":                                    "",
		"product_instance_type_family":                             "",
		"product_intel_avx2_available":                             "",
		"product_intel_avx2available":                              "",
		"product_intel_avx_available":                              "",
		"product_intel_turbo_available":                            "",
		"product_io":                                               "",
		"product_license_model":                                    "",
		"product_location":                                         "",
		"product_location_type":                                    "",
		"product_logs_destination":                                 "",
		"product_marketoption":                                     "",
		"product_maximum_extended_storage":                         "",
		"product_max_iops_burst_performance":                       "",
		"product_max_iopsvolume":                                   "",
		"product_max_throughputvolume":                             "",
		"product_max_volume_size":                                  "",
		"product_memory":                                           "",
		"product_memory_gib":                                       "",
		"product_memorytype":                                       "",
		"product_message_delivery_frequency":                       "",
		"product_message_delivery_order":                           "",
		"product_metering_type":                                    "",
		"product_min_volume_size":                                  "",
		"product_network_performance":                              "",
		"product_newcode":                                          "",
		"product_normalization_size_factor":                        "",
		"product_operating_system":                                 "",
		"product_operation":                                        "",
		"product_origin":                                           "",
		"product_overage_type":                                     "",
		"product_overhead":                                         "",
		"product_pack_size":                                        "",
		"product_parameter_type":                                   "",
		"product_physical_cpu":                                     "",
		"product_physical_gpu":                                     "",
		"product_physical_processor":                               "",
		"product_platoclassificationtype":                          "",
		"product_platoinstancename":                                "",
		"product_platoinstancetype":                                "",
		"product_platopricingtype":                                 "",
		"product_pre_installed_sw":                                 "",
		"product_pricingplan":                                      "",
		"product_pricing_unit":                                     "",
		"product_processor_architecture":                           "",
		"product_processor_features":                               "",
		"product_product_family":                                   "",
		"product_product_name":                                     "",
		"product_protocol":                                         "",
		"product_provider":                                         "",
		"product_provisioned":                                      "",
		"product_q_present":                                        "",
		"product_queue_type":                                       "",
		"product_recipient":                                        "",
		"product_region":                                           "",
		"product_region_code":                                      "",
		"product_request_description":                              "",
		"product_request_type":                                     "",
		"product_resource_endpoint":                                "",
		"product_resource_price_group":                             "",
		"product_routing_target":                                   "",
		"product_routing_type":                                     "",
		"product_scan_type":                                        "",
		"product_servicecode":                                      "",
		"product_servicename":                                      "",
		"product_size_flex":                                        "",
		"product_sku":                                              "",
		"product_software_type":                                    "",
		"product_standard_group":                                   "",
		"product_standard_storage":                                 "",
		"product_standard_storage_retention_included":              "",
		"product_steps":                                            "",
		"product_storage":                                          "",
		"product_storage_capacity_quota":                           "",
		"product_storage_class":                                    "",
		"product_storage_family":                                   "",
		"product_storage_media":                                    "",
		"product_storage_tier":                                     "",
		"product_storage_type":                                     "",
		"product_subcategory":                                      "",
		"product_subscription_type":                                "",
		"product_subservice":                                       "",
		"product_tenancy":                                          "",
		"product_throughput":                                       "",
		"product_throughput_capacity":                              "",
		"product_tiertype":                                         "",
		"product_to_location":                                      "",
		"product_to_location_type":                                 "",
		"product_to_region_code":                                   "",
		"product_transfer_type":                                    "",
		"product_type":                                             "",
		"product_usage_family":                                     "",
		"product_usage_group":                                      "",
		"product_usagetype":                                        "",
		"product_usage_volume":                                     "",
		"product_vaulttype":                                        "",
		"product_vcpu":                                             "",
		"product_version":                                          "",
		"product_volume_api_name":                                  "",
		"product_volume_type":                                      "",
		"product_vpcnetworkingsupport":                             "",
		"product_with_active_users":                                "",
		"reservation":                                              "A map column for where each key-value pair is an additional reservation attribute and its value.",
		"reservation_amortized_upfront_cost_for_usage":             "The initial upfront payment for all upfront RIs and partial upfront RIs amortized for usage time. The value is equal to: RIAmortizedUpfrontFeeForBillingPeriod * The normalized usage amount for DiscountedUsage line items / The normalized usage amount for the RIFee. Because there are no upfront payments for no upfront RIs, the value for a no upfront RI is 0. We do not provide this value for Dedicated Host reservations at this time. The change will be made in a future update.",
		"reservation_amortized_upfront_fee_for_billing_period":     "Describes how much of the upfront fee for this reservation is costing you for the billing period. The initial upfront payment for all upfront RIs and partial upfront RIs, amortized over this month. Because there are no upfront fees for no upfront RIs, the value for no upfront RIs is 0. We do not provide this value for Dedicated Host reservations at this time. The change will be made in a future update.",
		"reservation_reservation_arn":                              "The Amazon Resource Name (ARN) of the RI that this line item benefited from. This is also called the 'RI Lease ID'. This is a unique identifier of this particular AWS Reserved Instance. The value string also contains the AWS service name and the Region where the RI was purchased.",
		"reservation_availability_zone":                            "The Availability Zone of the resource that is associated with this line item.",
		"reservation_effective_cost":                               "The sum of both the upfront and hourly rate of your RI, averaged into an effective hourly rate. EffectiveCost is calculated by taking the amortizedUpfrontCostForUsage and adding it to the recurringFeeForUsage.",
		"reservation_end_time":                                     "The end date of the associated RI lease term.",
		"reservation_modification_status":                          "Shows whether the RI lease was modified or if it is unaltered.",
		"reservation_net_amortized_upfront_cost_for_usage":         "The initial upfront payment for All Upfront RIs and Partial Upfront RIs amortized for usage time, if applicable.",
		"reservation_net_amortized_upfront_fee_for_billing_period": "The cost of the reservation's upfront fee for the billing period.",
		"reservation_net_effective_cost":                           "The sum of both the upfront fee and the hourly rate of your RI, averaged into an effective hourly rate.",
		"reservation_net_recurring_fee_for_usage":                  "The after-discount cost of the recurring usage fee.",
		"reservation_net_unused_amortized_upfront_fee_for_billing_period":  "The net unused amortized upfront fee for the billing period.",
		"reservation_net_unused_recurring_fee":                             "The recurring fees associated with unused reservation hours for Partial Upfront and No Upfront RIs after discounts.",
		"reservation_net_upfront_value":                                    "The upfront value of the RI with discounts applied.",
		"reservation_normalized_units_per_reservation":                     "The number of normalized units for each instance of a reservation subscription.",
		"reservation_number_of_reservations":                               "The number of reservations that are covered by this subscription. For example, one RI subscription might have four associated RI reservations.",
		"reservation_recurring_fee_for_usage":                              "The recurring fee amortized for usage time, for partial upfront RIs and no upfront RIs. The value is equal to: The unblended cost of the RIFee * The sum of the normalized usage amount of Usage line items / The normalized usage amount of the RIFee for size flexible Reserved Instances. Because all upfront RIs don't have recurring fee payments greater than 0, the value for all upfront RIs is 0.",
		"reservation_start_time":                                           "The start date of the term of the associated Reserved Instance.",
		"reservation_subscription_id":                                      "A unique identifier that maps a line item with the associated offer. We recommend you use the RI ARN as your identifier of an AWS Reserved Instance, but both can be used.",
		"reservation_total_reserved_normalized_units":                      "The total number of reserved normalized units for all instances for a reservation subscription. AWS computes total normalized units by multiplying the reservation_normalized_units_per_reservation with reservation_number_of_reservations.",
		"reservation_total_reserved_units":                                 "TotalReservedUnits populates for both Fee and RIFee line items with distinct values. Fee line items: The total number of units reserved, for the total quantity of leases purchased in your subscription for the entire term. This is calculated by multiplying the NumberOfReservations with UnitsPerReservation. For example, 5 RIs x 744 hours per month x 12 months = 44,640. RIFee line items (monthly recurring costs): The total number of available units in your subscription, such as the total number of Amazon EC2 hours in a specific RI subscription. For example, 5 RIs x 744 hours = 3,720.",
		"reservation_units_per_reservation":                                "UnitsPerReservation populates for both Fee and RIFee line items with distinct values. Fee line items: The total number of units reserved for the subscription, such as the total number of RI hours purchased for the term of the subscription. For example 744 hours per month x 12 months = 8,928 total hours/units. RIFee line items (monthly recurring costs): The total number of available units in your subscription, such as the total number of Amazon EC2 hours in a specific RI subscription. For example, 1 unit x 744 hours = 744.",
		"reservation_unused_amortized_upfront_fee_for_billing_period":      "The amortized-upfront-fee-for-billing-period-column amortized portion of the initial upfront fee for all upfront RIs and partial upfront RIs. Because there are no upfront payments for no upfront RIs, the value for no upfront RIs is 0. We do not provide this value for Dedicated Host reservations at this time. The change will be made in a future update.",
		"reservation_unused_normalized_unit_quantity":                      "The number of unused normalized units for a size-flexible Regional RI that you didn't use during this billing period.",
		"reservation_unused_quantity":                                      "The number of RI hours that you didn't use during this billing period.",
		"reservation_unused_recurring_fee":                                 "The recurring fees associated with your unused reservation hours for partial upfront and no upfront RIs. Because all upfront RIs don't have recurring fees greater than 0, the value for All Upfront RIs is 0.",
		"reservation_upfront_value":                                        "The upfront price paid for your AWS Reserved Instance. For no upfront RIs, this value is 0.",
		"resource_tags":                                                    "A map where each entry is a resource tag key-value pair. This can be used to find information about the specific resources covered by a line item.",
		"resource_tags_aws_cloudformation_stack_name":                      "",
		"resource_tags_aws_created_by":                                     "",
		"resource_tags_user_cost_center":                                   "",
		"resource_tags_user_created_by_actor":                              "",
		"resource_tags_user_department":                                    "",
		"resource_tags_user_name":                                          "",
		"resource_tags_user_owner":                                         "",
		"savings_plan_amortized_upfront_commitment_for_billing_period":     "The amount of upfront fee a Savings Plan subscription is costing you for the billing period. The initial upfront payment for All Upfront Savings Plan and Partial Upfront Savings Plan amortized over the current month. For No Upfront Savings Plan, the value is 0.",
		"savings_plan_end_time":                                            "The expiration date for the Savings Plan agreement.",
		"savings_plan_instance_type_family":                                "The instance family that is associated with the specified usage.",
		"savings_plan_net_amortized_upfront_commitment_for_billing_period": "The cost of a Savings Plan subscription upfront fee for the billing period.",
		"savings_plan_net_recurring_commitment_for_billing_period":         "The net unblended cost of the Savings Plan fee.",
		"savings_plan_net_savings_plan_effective_cost":                     "The effective cost for Savings Plans, which is your usage divided by the fees.",
		"savings_plan_offering_type":                                       "Describes the type of Savings Plan purchased.",
		"savings_plan_payment_option":                                      "The payment options available for your Savings Plan.",
		"savings_plan_purchase_term":                                       "Describes the duration, or term, of the Savings Plan.",
		"savings_plan_recurring_commitment_for_billing_period":             "The monthly recurring fee for your Savings Plan subscriptions. For example, the recurring monthly fee for a Partial Upfront Savings Plan or No Upfront Savings Plan.",
		"savings_plan_region":                                              "The AWS Region (geographic area) that hosts your AWS services. You can use this field to analyze spend across a particular AWS Region.",
		"savings_plan_savings_plan_arn":                                    "The unique Savings Plan identifier.",
		"savings_plan_savings_plan_effective_cost":                         "The proportion of the Savings Plan monthly commitment amount (upfront and recurring) that is allocated to each usage line.",
		"savings_plan_savings_plan_rate":                                   "The Savings Plan rate for the usage.",
		"savings_plan_start_time":                                          "The start date of the Savings Plan agreement.",
		"savings_plan_total_commitment_to_date":                            "The total amortized upfront commitment and recurring commitment to date, for that hour.",
		"savings_plan_used_commitment":                                     "The total dollar amount of the Savings Plan commitment used. (SavingsPlanRate multiplied by usage).",
		"split_line_item_actual_usage":                                     "The usage for vCPU or memory (based on line_item_usage_type) you incurred for the specified time period for the Amazon ECS task.",
		"split_line_item_net_split_cost":                                   "The effective cost for Amazon ECS tasks after all discounts have been applied.",
		"split_line_item_net_unused_cost":                                  "The effective unused cost for Amazon ECS tasks after all discounts have been applied.",
		"split_line_item_parent_resource_id":                               "The resource ID of the parent EC2 instance associated with the Amazon ECS task (referenced in the line_item_resourceId column). The parent resource ID implies that the ECS task workload for the specified time period ran on the parent EC2 instance. This applies only for Amazon ECS tasks with EC2 launch type.",
		"split_line_item_public_on_demand_split_cost":                      "The cost for vCPU or memory (based on line_item_usage_type) allocated for the time period to the Amazon ECS task based on public On-Demand Instance rates (referenced in the pricing_public_on_demand_rate column).",
		"split_line_item_public_on_demand_unused_cost":                     "The unused cost for vCPU or memory (based on line_item_usage_type) allocated for the time period to the Amazon ECS task based on public On-Demand Instance rates. Unused costs are costs associated with resources (CPU or memory) on the EC2 instance (referenced in the split_line_item_parent_resource_id column) that were not utilized for the specified time period.",
		"split_line_item_reserved_usage":                                   "The usage for vCPU or memory (based on line_item_usage_type) that you configured for the specified time period for the Amazon ECS task.",
		"split_line_item_split_cost":                                       "The cost for vCPU or memory (based on line_item_usage_type) allocated for the time period to the Amazon ECS task. This includes amortized costs if the EC2 instance (referenced in the split_line_item_parent_resource_id column) has upfront or partial upfront charges for reservations or Savings Plans.",
		"split_line_item_split_usage":                                      "The usage for vCPU or memory (based on line_item_usage_type) allocated for the specified time period to the Amazon ECS task. This is defined as the maximum usage of split_line_item_reserved_usage or split_line_item_actual_usage.",
		"split_line_item_split_usage_ratio":                                "The ratio of vCPU or memory (based on line_item_usage_type) allocated to the Amazon ECS task compared to the overall CPU or memory available on the EC2 instance (referenced in the split_line_item_parent_resource_id column).",
		"split_line_item_unused_cost":                                      "The unused cost for vCPU or memory (based on line_item_usage_type) allocated for the time period to the Amazon ECS task. Unused costs are costs associated with resources (CPU or memory) on the EC2 instance (referenced in the split_line_item_parent_resource_id column) that were not utilized for the specified time period. This includes amortized costs if the EC2 instance (split_line_item_parent_resource_id) has upfront or partial upfront charges for reservations or Savings Plans.",
	}
}
