package cost_and_usage_focus

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"log/slog"
	"strconv"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/error_types"
)

type CostAndUsageFocusMapper struct {
}

func NewCostAndUsageFocusMapper() *CostAndUsageFocusMapper {
	return &CostAndUsageFocusMapper{}
}

func (m *CostAndUsageFocusMapper) Identifier() string {
	return "cost_and_usage_focus_mapper"
}

func (m *CostAndUsageFocusMapper) Map(_ context.Context, a any, opts ...mappers.MapOption) (*CostUsageFocus, error) {
	var input []byte

	// apply opts
	var config = &mappers.MapConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}

	// validate input type
	switch v := a.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	default:
		slog.Error("CostAndUsageFocusMapper.Map failed to map row  due to invalid type", "expected", "[]byte or string", "got", v)
		return nil, error_types.NewRowErrorWithMessage("unable to map row, invalid type received")
	}

	// read CSV line
	reader := csv.NewReader(bytes.NewReader(input))
	record, err := reader.Read()
	if err != nil {
		slog.Error("CostAndUsageFocusMapper.Map failed to read CSV line", "error", err)
		return nil, error_types.NewRowErrorWithMessage("failed to read log line")
	}

	// validate header/value count
	if len(record) != len(config.Header) {
		slog.Error("CostAndUsageFocusMapper.Map failed to map row due to header/value count mismatch", "expected", len(config.Header), "got", len(record))
		return nil, error_types.NewRowErrorWithMessage("row fields doesn't match count of headers")
	}

	// create a new CostUsageFocus object, with initialized maps
	output := NewCostUsageFocus()

	// map to struct
	for i, value := range record {
		field := config.Header[i]

		switch field {
		case "AvailabilityZone":
			output.AvailabilityZone = &value
		case "BilledCost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.BilledCost = &f
			}
		case "BillingAccountId":
			output.BillingAccountId = &value
		case "BillingAccountName":
			output.BillingAccountName = &value
		case "BillingCurrency":
			output.BillingCurrency = &value
		case "BillingPeriodEnd":
			if t, err := helpers.ParseTime(value); err == nil {
				output.BillingPeriodEnd = &t
			}
		case "BillingPeriodStart":
			if t, err := helpers.ParseTime(value); err == nil {
				output.BillingPeriodStart = &t
			}
		case "ChargeCategory":
			output.ChargeCategory = &value
		case "ChargeClass":
			output.ChargeClass = &value
		case "ChargeDescription":
			output.ChargeDescription = &value
		case "ChargeFrequency":
			output.ChargeFrequency = &value
		case "ChargePeriodEnd":
			if t, err := helpers.ParseTime(value); err == nil {
				output.ChargePeriodEnd = &t
			}
		case "ChargePeriodStart":
			if t, err := helpers.ParseTime(value); err == nil {
				output.ChargePeriodStart = &t
			}
		case "CommitmentDiscountCategory":
			output.CommitmentDiscountCategory = &value
		case "CommitmentDiscountId":
			output.CommitmentDiscountId = &value
		case "CommitmentDiscountName":
			output.CommitmentDiscountName = &value
		case "CommitmentDiscountStatus":
			output.CommitmentDiscountStatus = &value
		case "CommitmentDiscountType":
			output.CommitmentDiscountType = &value
		case "ConsumedQuantity":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ConsumedQuantity = &f
			}
		case "ConsumedUnit":
			output.ConsumedUnit = &value
		case "ContractedCost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ContractedCost = &f
			}
		case "ContractedUnitPrice":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ContractedUnitPrice = &f
			}
		case "EffectiveCost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EffectiveCost = &f
			}
		case "InvoiceIssuerName":
			output.InvoiceIssuerName = &value
		case "ListCost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ListCost = &f
			}
		case "ListUnitPrice":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ListUnitPrice = &f
			}
		case "PricingCategory":
			output.PricingCategory = &value
		case "PricingQuantity":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.PricingQuantity = &f
			}
		case "PricingUnit":
			output.PricingUnit = &value
		case "ProviderName":
			output.ProviderName = &value
		case "PublisherName":
			output.PublisherName = &value
		case "RegionId":
			output.RegionId = &value
		case "RegionName":
			output.RegionName = &value
		case "ResourceId":
			output.ResourceId = &value
		case "ResourceName":
			output.ResourceName = &value
		case "ResourceType":
			output.ResourceType = &value
		case "ServiceCategory":
			output.ServiceCategory = &value
		case "ServiceName":
			output.ServiceName = &value
		case "SkuId":
			output.SkuId = &value
		case "SkuPriceId":
			output.SkuPriceId = &value
		case "SubAccountId":
			output.SubAccountId = &value
		case "SubAccountName":
			output.SubAccountName = &value
		case "Tags":
			tags := make(map[string]string)
			if err := json.Unmarshal([]byte(value), &tags); err == nil {
				output.Tags = &tags
			}
		case "x_CostCategories":
			costCategories := make(map[string]string)
			if err := json.Unmarshal([]byte(value), &costCategories); err == nil {
				output.XCostCategories = &costCategories
			}
		case "x_Discounts":
			discounts := make(map[string]string)
			if err := json.Unmarshal([]byte(value), &discounts); err == nil {
				output.XDiscounts = &discounts
			}
		case "x_Operation":
			output.XOperation = &value
		case "x_ServiceCode":
			output.XServiceCode = &value
		case "x_UsageType":
			output.XUsageType = &value
		}
	}

	return output, nil
}
