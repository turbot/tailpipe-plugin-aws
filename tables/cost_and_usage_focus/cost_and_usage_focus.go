package cost_and_usage_focus

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type CostUsageFocus struct {
	schema.CommonFields

	AvailabilityZone           *string            `json:"AvailabilityZone,omitempty" parquet:"name=availability_zone"`
	BilledCost                 *float64           `json:"BilledCost,omitempty" parquet:"name=billed_cost"`
	BillingAccountId           *string            `json:"BillingAccountId,omitempty" parquet:"name=billing_account_id"`
	BillingAccountName         *string            `json:"BillingAccountName,omitempty" parquet:"name=billing_account_name"`
	BillingCurrency            *string            `json:"BillingCurrency,omitempty" parquet:"name=billing_currency"`
	BillingPeriodEnd           *time.Time         `json:"BillingPeriodEnd,omitempty" parquet:"name=billing_period_end"`
	BillingPeriodStart         *time.Time         `json:"BillingPeriodStart,omitempty" parquet:"name=billing_period_start"`
	ChargeCategory             *string            `json:"ChargeCategory,omitempty" parquet:"name=charge_category"`
	ChargeClass                *string            `json:"ChargeClass,omitempty" parquet:"name=charge_class"`
	ChargeDescription          *string            `json:"ChargeDescription,omitempty" parquet:"name=charge_description"`
	ChargeFrequency            *string            `json:"ChargeFrequency,omitempty" parquet:"name=charge_frequency"`
	ChargePeriodEnd            *time.Time         `json:"ChargePeriodEnd,omitempty" parquet:"name=charge_period_end"`
	ChargePeriodStart          *time.Time         `json:"ChargePeriodStart,omitempty" parquet:"name=charge_period_start"`
	CommitmentDiscountCategory *string            `json:"CommitmentDiscountCategory,omitempty" parquet:"name=commitment_discount_category"`
	CommitmentDiscountId       *string            `json:"CommitmentDiscountId,omitempty" parquet:"name=commitment_discount_id"`
	CommitmentDiscountName     *string            `json:"CommitmentDiscountName,omitempty" parquet:"name=commitment_discount_name"`
	CommitmentDiscountStatus   *string            `json:"CommitmentDiscountStatus,omitempty" parquet:"name=commitment_discount_status"`
	CommitmentDiscountType     *string            `json:"CommitmentDiscountType,omitempty" parquet:"name=commitment_discount_type"`
	ConsumedQuantity           *float64           `json:"ConsumedQuantity,omitempty" parquet:"name=consumed_quantity"`
	ConsumedUnit               *string            `json:"ConsumedUnit,omitempty" parquet:"name=consumed_unit"`
	ContractedCost             *float64           `json:"ContractedCost,omitempty" parquet:"name=contracted_cost"`
	ContractedUnitPrice        *float64           `json:"ContractedUnitPrice,omitempty" parquet:"name=contracted_unit_price"`
	EffectiveCost              *float64           `json:"EffectiveCost,omitempty" parquet:"name=effective_cost"`
	InvoiceIssuerName          *string            `json:"InvoiceIssuerName,omitempty" parquet:"name=invoice_issuer_name"`
	ListCost                   *float64           `json:"ListCost,omitempty" parquet:"name=list_cost"`
	ListUnitPrice              *float64           `json:"ListUnitPrice,omitempty" parquet:"name=list_unit_price"`
	PricingCategory            *string            `json:"PricingCategory,omitempty" parquet:"name=pricing_category"`
	PricingQuantity            *float64           `json:"PricingQuantity,omitempty" parquet:"name=pricing_quantity"`
	PricingUnit                *string            `json:"PricingUnit,omitempty" parquet:"name=pricing_unit"`
	ProviderName               *string            `json:"ProviderName,omitempty" parquet:"name=provider_name"`
	PublisherName              *string            `json:"PublisherName,omitempty" parquet:"name=publisher_name"`
	RegionId                   *string            `json:"RegionId,omitempty" parquet:"name=region_id"`
	RegionName                 *string            `json:"RegionName,omitempty" parquet:"name=region_name"`
	ResourceId                 *string            `json:"ResourceId,omitempty" parquet:"name=resource_id"`
	ResourceName               *string            `json:"ResourceName,omitempty" parquet:"name=resource_name"`
	ResourceType               *string            `json:"ResourceType,omitempty" parquet:"name=resource_type"`
	ServiceCategory            *string            `json:"ServiceCategory,omitempty" parquet:"name=service_category"`
	ServiceName                *string            `json:"ServiceName,omitempty" parquet:"name=service_name"`
	SkuId                      *string            `json:"SkuId,omitempty" parquet:"name=sku_id"`
	SkuPriceId                 *string            `json:"SkuPriceId,omitempty" parquet:"name=sku_price_id"`
	SubAccountId               *string            `json:"SubAccountId,omitempty" parquet:"name=sub_account_id"`
	SubAccountName             *string            `json:"SubAccountName,omitempty" parquet:"name=sub_account_name"`
	Tags                       *map[string]string `json:"Tags,omitempty" parquet:"name=tags"`                          // -- MAP
	XCostCategories            *map[string]string `json:"x_CostCategories,omitempty" parquet:"name=x_cost_categories"` // -- MAP
	XDiscounts                 *map[string]string `json:"x_Discounts,omitempty" parquet:"name=x_discounts"`            // -- MAP
	XOperation                 *string            `json:"x_Operation,omitempty" parquet:"name=x_operation"`
	XServiceCode               *string            `json:"x_ServiceCode,omitempty" parquet:"name=x_service_code"`
	XUsageType                 *string            `json:"x_UsageType,omitempty" parquet:"name=x_usage_type"`
}

func (c *CostUsageFocus) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"availability_zone":            "A provider-assigned identifier for a physically separated and isolated area within a Region that provides high availability and fault tolerance.",
		"billed_cost":                  "A charge that is the basis for invoicing, inclusive of all reduced rates and discounts while excluding the amortization of relevant purchases paid to cover future eligible charges.",
		"billing_account_id":           "A provider-assigned identifier for a billing account.",
		"billing_account_name":         "A provider-assigned name for a billing account.",
		"billing_currency":             "An identifier that represents the currency that a charge for resources or services was billed in.",
		"billing_period_end":           "The end date and time of the billing period.",
		"billing_period_start":         "The start date and time of the billing period.",
		"charge_category":              "An indicator of whether the row represents an upfront or recurring fee, cost of usage that already occurred, an after-the-fact adjustment (for example, credits), or taxes.",
		"charge_class":                 "An indicator of whether the row represents a regular charge, or a correction to one or more previous charges.",
		"charge_description":           "A high-level context of a row without requiring additional discovery.",
		"charge_frequency":             "An indicator of how often a charge will occur.",
		"charge_period_end":            "The end date and time of the charge period.",
		"charge_period_start":          "The start date and time of the charge period.",
		"commitment_discount_category": "An indicator of whether the commitment-based discount identified in the commitment_discount_id column is based on usage quantity or cost (that is, spend).",
		"commitment_discount_id":       "A provider-assigned identifier for a commitment-based discount.",
		"commitment_discount_name":     "The display name assigned to a commitment-based discount.",
		"commitment_discount_status":   "An indicator of whether the charge corresponds to a used or unused commitment discount.",
		"commitment_discount_type":     "A provider-assigned name to identify the type of commitment-based discount applied to the row.",
		"consumed_quantity":            "The volume of a given resource or service used or purchased based on the consumed_unit.",
		"consumed_unit":                "A provider-assigned measurement unit indicating how a provider measures usage of a given SKU associated with a resource or service.",
		"contracted_cost":              "The cost calculated by multiplying contracted_unit_price and the corresponding pricing_quantity.",
		"contracted_unit_price":        "The agreed-upon unit price for a single pricing_unit of the associated SKU, inclusive of any negotiated discounts while excluding negotiated commitment-based discounts or any other discounts.",
		"effective_cost":               "A cost that includes all reduced rates and discounts, augmented with the amortization of relevant purchases (one-time or recurring) paid to cover future eligible charges.",
		"invoice_issuer_name":          "An entity responsible for invoicing the sources or services consumed. It is commonly used for cost analysis and reporting scenarios.",
		"list_cost":                    "The cost calculated by multiplying list_unit_price and the corresponding pricing_quantity.",
		"list_unit_price":              "The suggested unit price, published by the provider, for a single pricing_unit of the associated SKU, excluding any discounts.",
		"pricing_category":             "The pricing model used for a charge at the time of use or purchase.",
		"pricing_quantity":             "The volume of a given SKU associated with a resource or service used or purchased, based on the pricing_unit.",
		"pricing_unit":                 "A provider-assigned measurement unit for determining unit prices, indicating how the provider rates measured usage and purchase quantities after applying pricing rules such as block pricing.",
		"provider_name":                "The entity that made the resources or services available for purchase.",
		"publisher_name":               "The entity that produced the resources or services that were purchased.",
		"region_id":                    "A provider-assigned identifier for an isolated geographic area where a resource is provisioned or a service is provided.",
		"region_name":                  "The name of an isolated geographic area where a resource is provisioned or a service is provided.",
		"resource_id":                  "A provider-assigned identifier for a resource.",
		"resource_name":                "A display name assigned to a resource.",
		"resource_type":                "The type of resource the charge applies to.",
		"service_category":             "The highest-level classification of a service based on the core function of the service.",
		"service_name":                 "A display name for the offering that was purchased.",
		"sku_id":                       "A unique identifier that defines a provider-supported construct for organizing properties that are common across one or more SKU prices.",
		"sku_price_id":                 "A unique identifier that defines the unit price used to calculate the charge.",
		"sub_account_id":               "An ID assigned to a grouping of resources or services, often used to manage access and/or cost.",
		"sub_account_name":             "A name assigned to a grouping of resources or services, often used to manage access and/or cost.",
		"tags":                         "The set of tags assigned to tag sources that also account for potential provider-defined or user-defined tag evaluations.",
		"x_cost_categories":            "A map column containing key-value pairs of the cost categories and their values for a given line item.",
		"x_discounts":                  "A map column containing key-value pairs of any specific discounts that apply to this line item.",
		"x_operation":                  "The specific AWS operation covered by this line item.",
		"x_service_code":               "The code of the service used in this line item.",
		"x_usage_type":                 "The usage details of the line item.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "The list of ARNs associated with a cost and usage report. If ResourceId starts with 'arn:', it is included in this list.",
		"tp_index":     "The account ID associated with the report, determined based on the following priority: SubAccountId, BillingAccountId, or a default value if neither is available.",
		"tp_timestamp": "The timestamp representing the start or end date of the usage. If available, ChargePeriodStart is used first, followed by ChargePeriodEnd, BillingPeriodStart, or BillingPeriodEnd in that order. The timestamp is stored in ISO 8601 format.",
	}
}
