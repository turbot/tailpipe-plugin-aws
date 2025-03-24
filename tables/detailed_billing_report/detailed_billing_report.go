package detailed_billing_report

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type DetailedBillingReport struct {
	schema.CommonFields

	AvailabilityZone   *string            `json:"AvailabilityZone,omitempty" parquet:"name=availability_zone"`
	BillingPeriodEnd   *time.Time         `json:"BillingPeriodEndDate,omitempty" parquet:"name=billing_period_end"`
	BillingPeriodStart *time.Time         `json:"BillingPeriodStartDate,omitempty" parquet:"name=billing_period_start"`
	BlendedCost        *float64           `json:"BlendedCost,omitempty" parquet:"name=blended_cost"`
	BlendedRate        *float64           `json:"BlendedRate,omitempty" parquet:"name=blended_rate"`
	CostBeforeTax      *float64           `json:"CostBeforeTax,omitempty" parquet:"name=cost_before_tax"`
	Credits            *float64           `json:"Credits,omitempty" parquet:"name=credits"`
	CurrencyCode       *string            `json:"CurrencyCode,omitempty" parquet:"name=currency_code"`
	InvoiceDate        *time.Time         `json:"InvoiceDate,omitempty" parquet:"name=invoice_date"`
	InvoiceID          *string            `json:"InvoiceID,omitempty" parquet:"name=invoice_id"`
	ItemDescription    *string            `json:"ItemDescription,omitempty" parquet:"name=item_description"`
	LinkedAccountId    *string            `json:"LinkedAccountId,omitempty" parquet:"name=linked_account_id"`
	LinkedAccountName  *string            `json:"LinkedAccountName,omitempty" parquet:"name=linked_account_name"`
	Operation          *string            `json:"Operation,omitempty" parquet:"name=operation"`
	PayerAccountId     *string            `json:"PayerAccountId,omitempty" parquet:"name=payer_account_id"`
	PayerAccountName   *string            `json:"PayerAccountName,omitempty" parquet:"name=payer_account_name"`
	PayerPONumber      *string            `json:"PayerPONumber,omitempty" parquet:"name=payer_po_number"`
	PricingPlanId      *string            `json:"PricingPlanId,omitempty" parquet:"name=pricing_plan_id"`
	ProductCode        *string            `json:"ProductCode,omitempty" parquet:"name=product_code"`
	ProductName        *string            `json:"ProductName,omitempty" parquet:"name=product_name"`
	RateId             *string            `json:"RateId,omitempty" parquet:"name=rate_id"`
	RecordID           *string            `json:"RecordID,omitempty" parquet:"name=record_id"`
	RecordType         *string            `json:"RecordType,omitempty" parquet:"name=record_type"`
	ReservedInstance   *string            `json:"ReservedInstance,omitempty" parquet:"name=reserved_instance"`
	ResourceId         *string            `json:"ResourceId,omitempty" parquet:"name=resource_id"`
	SellerOfRecord     *string            `json:"SellerOfRecord,omitempty" parquet:"name=seller_of_record"`
	SubscriptionId     *string            `json:"SubscriptionId,omitempty" parquet:"name=subscription_id"`
	TaxAmount          *float64           `json:"TaxAmount,omitempty" parquet:"name=tax_amount"`
	TaxationAddress    *string            `json:"TaxationAddress,omitempty" parquet:"name=taxation_address"`
	TaxType            *string            `json:"TaxType,omitempty" parquet:"name=tax_type"`
	TotalCost          *float64           `json:"TotalCost,omitempty" parquet:"name=total_cost"`
	UnBlendedCost      *float64           `json:"UnBlendedCost,omitempty" parquet:"name=un_blended_cost"`
	UnBlendedRate      *float64           `json:"UnBlendedRate,omitempty" parquet:"name=un_blended_rate"`
	UsageEndDate       *time.Time         `json:"UsageEndDate,omitempty" parquet:"name=usage_end_date"`
	UsageQuantity      *float64           `json:"UsageQuantity,omitempty" parquet:"name=usage_quantity"`
	UsageStartDate     *time.Time         `json:"UsageStartDate,omitempty" parquet:"name=usage_start_date"`
	UsageType          *string            `json:"UsageType,omitempty" parquet:"name=usage_type"`
	Tags               *map[string]string `json:"Tags,omitempty" parquet:"name=tags"`
}

func (c *DetailedBillingReport) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"availability_zone":    "Availability zone where the resource was used.",
		"billing_period_end":   "End date of the billing period.",
		"billing_period_start": "Start date of the billing period.",
		"blended_cost":         "Blended cost based on usage and rate.",
		"blended_rate":         "Blended rate applied.",
		"cost_before_tax":      "Total cost before taxes.",
		"credits":              "Credits applied.",
		"currency_code":        "Currency in which billing is done.",
		"invoice_date":         "Date the invoice was issued.",
		"invoice_id":           "Unique identifier for the invoice.",
		"item_description":     "Description of the billed item.",
		"linked_account_id":    "Linked account under the payer.",
		"linked_account_name":  "Name of the linked account.",
		"operation":            "Operation performed during usage.",
		"payer_account_id":     "AWS payer account ID.",
		"payer_account_name":   "Name of the payer account.",
		"payer_po_number":      "Purchase order number from payer.",
		"pricing_plan_id":      "ID of the pricing plan associated with usage.",
		"product_code":         "AWS product code.",
		"product_name":         "Name of the AWS product.",
		"rate_id":              "ID representing the rate applied.",
		"record_id":            "Unique identifier for the billing record.",
		"record_type":          "Type of the billing record, e.g., LineItem or PayerLineItem.",
		"reserved_instance":    "Indicates whether reserved instance pricing is applied.",
		"resource_id":          "Identifier of the AWS resource.",
		"seller_of_record":     "Entity selling the product.",
		"subscription_id":      "Unique ID of the AWS subscription.",
		"tags":                 "A map where each entry is a resource tag key-value pair. This can be used to find information about the specific resources covered by a line item.",
		"tax_amount":           "Amount of tax applied.",
		"tax_type":             "Type of tax applied.",
		"taxation_address":     "Address used for tax calculation.",
		"total_cost":           "Total cost including tax and credits.",
		"un_blended_cost":      "Unblended cost without applying cost allocation.",
		"un_blended_rate":      "Unblended rate without applying cost allocation.",
		"usage_end_date":       "End date of usage.",
		"usage_quantity":       "Amount of usage.",
		"usage_start_date":     "Start date of usage.",
		"usage_type":           "Type of usage being billed.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "The list of ARNs associated with this billing record.",
		"tp_index":     "The account ID used for indexing the record.",
		"tp_timestamp": "Timestamp of the billing entry based on usage or billing period.",
	}
}
