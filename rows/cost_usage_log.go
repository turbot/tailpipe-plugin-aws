package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type CostAndUsageLog struct {
	schema.CommonFields

	InvoiceID              *string    `json:"invoice_id,omitempty"`
	PayerAccountId         *string    `json:"payer_account_id,omitempty"`
	LinkedAccountId        *string    `json:"linked_account_id,omitempty"`
	RecordType             *string    `json:"record_type,omitempty"`
	RecordID               *string    `json:"record_id,omitempty"`
	BillingPeriodStartDate *time.Time `json:"billing_period_start_date,omitempty"`
	BillingPeriodEndDate   *time.Time `json:"billing_period_end_date,omitempty"`
	InvoiceDate            *time.Time `json:"invoice_date,omitempty"`
	PayerAccountName       *string    `json:"payer_account_name,omitempty"`
	LinkedAccountName      *string    `json:"linked_account_name,omitempty"`
	TaxationAddress        *string    `json:"taxation_address,omitempty"`
	PayerPONumber          *string    `json:"payer_po_number,omitempty"`
	ProductCode            *string    `json:"product_code,omitempty"`
	ProductName            *string    `json:"product_name,omitempty"`
	SellerOfRecord         *string    `json:"seller_of_record,omitempty"`
	UsageType              *string    `json:"usage_type,omitempty"`
	Operation              *string    `json:"operation,omitempty"`
	RateId                 *string    `json:"rate_id,omitempty"`
	ItemDescription        *string    `json:"item_description,omitempty"`
	UsageStartDate         *time.Time `json:"usage_start_date,omitempty"`
	UsageEndDate           *time.Time `json:"usage_end_date,omitempty"`
	UsageQuantity          *float64   `json:"usage_quantity,omitempty"`
	BlendedRate            *string    `json:"blended_rate,omitempty"`
	CurrencyCode           *string    `json:"currency_code,omitempty"`
	CostBeforeTax          *float64   `json:"cost_before_tax,omitempty"`
	Credits                *float64   `json:"credits,omitempty"`
	TaxAmount              *float64   `json:"tax_amount,omitempty"`
	TaxType                *string    `json:"tax_type,omitempty"`
	TotalCost              *float64   `json:"total_cost,omitempty"`
}
