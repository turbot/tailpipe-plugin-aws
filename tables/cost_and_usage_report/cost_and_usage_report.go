package cost_and_usage_report

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// CUR 2.0 Schema: https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cur2.html
type CostAndUsageReport struct {
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
	LineItemUsageAmount                                      *float64                `json:"line_item_usage_amount,omitempty" parquet:"name=line_item_usage_amount"`
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
	ProductComment                                           *string                 `json:"product_comment,omitempty" parquet:"name=product_comment"`
	ProductFeeCode                                           *string                 `json:"product_fee_code,omitempty" parquet:"name=product_fee_code"`
	ProductFeeDescription                                    *string                 `json:"product_fee_description,omitempty" parquet:"name=product_fee_description"`
	ProductFromLocation                                      *string                 `json:"product_from_location,omitempty" parquet:"name=product_from_location"`
	ProductFromLocationType                                  *string                 `json:"product_from_location_type,omitempty" parquet:"name=product_from_location_type"`
	ProductFromRegionCode                                    *string                 `json:"product_from_region_code,omitempty" parquet:"name=product_from_region_code"`
	ProductInstanceSku                                       *string                 `json:"product_instancesku,omitempty" parquet:"name=product_instance_sku"`
	ProductInstanceFamily                                    *string                 `json:"product_instance_family,omitempty" parquet:"name=product_instance_family"`
	ProductInstanceType                                      *string                 `json:"product_instance_type,omitempty" parquet:"name=product_instance_type"`
	ProductLocation                                          *string                 `json:"product_location,omitempty" parquet:"name=product_location"`
	ProductLocationType                                      *string                 `json:"product_location_type,omitempty" parquet:"name=product_location_type"`
	ProductOperation                                         *string                 `json:"product_operation,omitempty" parquet:"name=product_operation"`
	ProductPricingUnit                                       *string                 `json:"product_pricing_unit,omitempty" parquet:"name=product_pricing_unit"`
	ProductProductFamily                                     *string                 `json:"product_product_family,omitempty" parquet:"name=product_product_family"`
	ProductRegionCode                                        *string                 `json:"product_region_code,omitempty" parquet:"name=product_region_code"`
	ProductSku                                               *string                 `json:"product_sku,omitempty" parquet:"name=product_sku"`
	ProductServiceCode                                       *string                 `json:"product_servicecode,omitempty" parquet:"name=product_service_code"`
	ProductToLocationType                                    *string                 `json:"product_to_location_type,omitempty" parquet:"name=product_to_location_type"`
	ProductToLocation                                        *string                 `json:"product_to_location,omitempty" parquet:"name=product_to_location"`
	ProductToRegionCode                                      *string                 `json:"product_to_region_code,omitempty" parquet:"name=product_to_region_code"`
	ProductUsageType                                         *string                 `json:"product_usagetype,omitempty" parquet:"name=product_usage_type"`
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

func (c *CostAndUsageReport) GetColumnDescriptions() map[string]string {
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
		"product_comment":                                          "A comment regarding the product.",
		"product_fee_code":                                         "Unique code for an AWS fee item.",
		"product_fee_description":                                  "Description of additional AWS fees (e.g., support, overage, licensing).",
		"product_from_location":                                    "Describes the location where the usage originated from.",
		"product_from_location_type":                               "Describes the location type where the usage originated from.",
		"product_from_region_code":                                 "Describes the source Region code for the AWS service.",
		"product_instance_sku":                                     "SKU identifier for a specific AWS instance type.",
		"product_instance_family":                                  "Describes your Amazon EC2 instance family. Amazon EC2 provides you with a large number of options across 10 different instance types, each with one or more size options, organized into distinct instance families optimized for different types of applications.",
		"product_instance_type":                                    "Describes the instance type, size, and family, which define the CPU, networking, and storage capacity of your instance.",
		"product_location":                                         "Describes the Region that your Amazon S3 bucket resides in.",
		"product_location_type":                                    "Describes the endpoint of your task.",
		"product_operation":                                        "Describes the specific AWS operation that this line item covers.",
		"product_pricing_unit":                                     "The smallest billing unit for an AWS service. For example, 0.01c per API call.",
		"product_product_family":                                   "The category for the type of product.",
		"product_region_code":                                      "A Region is a physical location around the world where data centers are clustered. AWS calls each group of logical data centers an Availability Zone (AZ). Each AWS Region consists of multiple, isolates, and physically separate AZs within a geographical area. The Region code attribute has the same name as an AWS Region, and specifies where the AWS service is available.",
		"product_sku":                                              "A unique code for a product. The SKU is created by combining the ProductCode, UsageType, and Operation. For size-flexible RIs, the SKU uses the instance that was used. For example, if you used a t2.micro instance and AWS applied a t2.small RI discount to the usage, the line item SKU is created with the t2.micro.",
		"product_service_code":                                     "This identifies the specific AWS service to the customer as a unique short abbreviation.",
		"product_to_location":                                      "Describes the location usage destination.",
		"product_to_location_type":                                 "Describes the destination location of the service usage.",
		"product_to_region_code":                                   "Describes the source Region code for the AWS service.",
		"product_usage_type":                                       "Describes the usage details of the line item.",
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
