package cost_and_usage_report

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/error_types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type CostAndUsageReportMapper struct {
	headers []string
}

func NewCostAndUsageReportMapper() *CostAndUsageReportMapper {
	return &CostAndUsageReportMapper{}
}

func (c *CostAndUsageReportMapper) Identifier() string {
	return "cost_and_usage_report_mapper"
}

func (c *CostAndUsageReportMapper) Map(_ context.Context, a any, opts ...mappers.MapOption[*CostUsageReport]) (*CostUsageReport, error) {
	var input []byte
	// apply opts
	for _, opt := range opts {
		opt(c)
	}

	// validate input type
	switch v := a.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	default:
		slog.Error("CostAndUsageReportMapper.Map failed to map row  due to invalid type", "expected", "[]byte or string", "got", v)
		return nil, error_types.NewRowErrorWithMessage("unable to map row, invalid type received")
	}

	// read CSV line
	reader := csv.NewReader(bytes.NewReader(input))
	record, err := reader.Read()
	if err != nil {
		slog.Error("CostAndUsageReportMapper.Map failed to read CSV line", "error", err)
		return nil, error_types.NewRowErrorWithMessage("failed to read log line")
	}

	// validate header/value count
	if len(record) != len(c.headers) {
		slog.Error("CostAndUsageReportMapper.Map failed to map row due to header/value count mismatch", "expected", len(c.headers), "got", len(record))
		return nil, error_types.NewRowErrorWithMessage("row fields doesn't match count of headers")
	}

	// create a new CostUsageReport object with initialised maps
	output := NewCostUsageReport()

	// map to struct (normalize headers)
	for i, value := range record {
		field := strings.ToLower(c.headers[i])

		switch field {
		case "bill_billing_entity":
			output.BillBillingEntity = &value
		case "bill_billing_period_start_date":
			if t, err := helpers.ParseTime(value); err == nil {
				output.BillBillingPeriodStartDate = &t
			}
		case "bill_billing_period_end_date":
			if t, err := helpers.ParseTime(value); err == nil {
				output.BillBillingPeriodEndDate = &t
			}
		case "bill_bill_type":
			output.BillBillType = &value
		case "bill_invoice_id":
			output.BillInvoiceId = &value
		case "bill_invoicing_entity":
			output.BillInvoicingEntity = &value
		case "bill_payer_account_id":
			output.BillPayerAccountId = &value
		case "bill_payer_account_name":
			output.BillPayerAccountName = &value
		case "cost_category":
			(*output.CostCategory)["cost_category"] = value
		case "discount":
			(*output.Discount)["discount"] = value
		case "discount_bundled_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.DiscountBundledDiscount = &f
				(*output.Discount)["bundled_discount"] = f
			}
		case "discount_total_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.DiscountTotalDiscount = &f
				(*output.Discount)["total_discount"] = f
			}
		case "identity_line_item_id":
			output.IdentityLineItemId = &value
		case "identity_time_interval":
			output.IdentityTimeInterval = &value
		case "line_item_availability_zone":
			output.LineItemAvailabilityZone = &value
		case "line_item_blended_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemBlendedCost = &f
			}
		case "line_item_blended_rate":
			output.LineItemBlendedRate = &value
		case "line_item_currency_code":
			output.LineItemCurrencyCode = &value
		case "line_item_legal_entity":
			output.LineItemLegalEntity = &value
		case "line_item_line_item_description":
			output.LineItemLineItemDescription = &value
		case "line_item_line_item_type":
			output.LineItemLineItemType = &value
		case "line_item_net_unblended_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemNetUnblendedCost = &f
			}
		case "line_item_net_unblended_rate":
			output.LineItemNetUnblendedRate = &value
		case "line_item_normalization_factor":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemNormalizationFactor = &f
			}
		case "line_item_normalized_usage_amount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemNormalizedUsageAmount = &f
			}
		case "line_item_operation":
			output.LineItemOperation = &value
		case "line_item_product_code":
			output.LineItemProductCode = &value
		case "line_item_resource_id":
			output.LineItemResourceId = &value
		case "line_item_tax_type":
			output.LineItemTaxType = &value
		case "line_item_unblended_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemUnblendedCost = &f
			}
		case "line_item_unblended_rate":
			output.LineItemUnblendedRate = &value
		case "line_item_usage_account_id":
			output.LineItemUsageAccountId = &value
		case "line_item_usage_account_name":
			output.LineItemUsageAccountName = &value
		case "line_item_usage_amount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.LineItemUsageAmount = &f
			}
		case "line_item_usage_end_date":
			if t, err := helpers.ParseTime(value); err == nil {
				output.LineItemUsageEndDate = &t
			}
		case "line_item_usage_start_date":
			if t, err := helpers.ParseTime(value); err == nil {
				output.LineItemUsageStartDate = &t
			}
		case "line_item_usage_type":
			output.LineItemUsageType = &value
		case "pricing_currency":
			output.PricingCurrency = &value
		case "pricing_lease_contract_length":
			output.PricingLeaseContractLength = &value
		case "pricing_offering_class":
			output.PricingOfferingClass = &value
		case "pricing_public_on_demand_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.PricingPublicOnDemandCost = &f
			}
		case "pricing_public_on_demand_rate":
			output.PricingPublicOnDemandRate = &value
		case "pricing_purchase_option":
			output.PricingPurchaseOption = &value
		case "pricing_rate_code":
			output.PricingRateCode = &value
		case "pricing_rate_id":
			output.PricingRateId = &value
		case "pricing_term":
			output.PricingTerm = &value
		case "pricing_unit":
			output.PricingUnit = &value
		case "product":
			(*output.Product)["product"] = value
		case "product_comment":
			output.ProductComment = &value
			(*output.Product)["comment"] = value
		case "product_fee_code":
			output.ProductFeeCode = &value
			(*output.Product)["fee_code"] = value
		case "product_fee_description":
			output.ProductFeeDescription = &value
			(*output.Product)["fee_description"] = value
		case "product_from_location":
			output.ProductFromLocation = &value
			(*output.Product)["from_location"] = value
		case "product_from_location_type":
			output.ProductFromLocationType = &value
			(*output.Product)["from_location_type"] = value
		case "product_from_region_code":
			output.ProductFromRegionCode = &value
			(*output.Product)["from_region_code"] = value
		case "product_instance_family":
			output.ProductInstanceFamily = &value
			(*output.Product)["instance_family"] = value
		case "product_instance_type":
			output.ProductInstanceType = &value
			(*output.Product)["instance_type"] = value
		case "product_instancesku":
			output.ProductInstanceSku = &value
			(*output.Product)["instancesku"] = value
		case "product_location":
			output.ProductLocation = &value
			(*output.Product)["location"] = value
		case "product_location_type":
			output.ProductLocationType = &value
			(*output.Product)["location_type"] = value
		case "product_operation":
			output.ProductOperation = &value
			(*output.Product)["operation"] = value
		case "product_pricing_unit":
			output.ProductPricingUnit = &value
			(*output.Product)["pricing_unit"] = value
		case "product_product_family":
			output.ProductProductFamily = &value
			(*output.Product)["product_family"] = value
		case "product_region_code":
			output.ProductRegionCode = &value
			(*output.Product)["region_code"] = value
		case "product_servicecode":
			output.ProductServiceCode = &value
			(*output.Product)["servicecode"] = value
		case "product_sku":
			output.ProductSku = &value
			(*output.Product)["sku"] = value
		case "product_to_location":
			output.ProductToLocation = &value
			(*output.Product)["to_location"] = value
		case "product_to_location_type":
			output.ProductToLocationType = &value
			(*output.Product)["to_location_type"] = value
		case "product_to_region_code":
			output.ProductToRegionCode = &value
			(*output.Product)["to_region_code"] = value
		case "product_usagetype":
			output.ProductUsageType = &value
			(*output.Product)["usagetype"] = value
		case "reservation_amortized_upfront_cost_for_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationAmortizedUpfrontCostForUsage = &f
				(*output.Reservation)["amortized_upfront_cost_for_usage"] = f
			}
		case "reservation_amortized_upfront_fee_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationAmortizedUpfrontFeeForBillingPeriod = &f
				(*output.Reservation)["amortized_upfront_fee_for_billing_period"] = f
			}
		case "reservation_availability_zone":
			output.ReservationAvailabilityZone = &value
			(*output.Reservation)["availability_zone"] = value
		case "reservation_effective_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationEffectiveCost = &f
				(*output.Reservation)["effective_cost"] = f
			}
		case "reservation_end_time":
			if t, err := helpers.ParseTime(value); err == nil {
				ts := t.Format(time.RFC3339)
				output.ReservationEndTime = &ts
				(*output.Reservation)["end_time"] = ts
			}
		case "reservation_modification_status":
			output.ReservationModificationStatus = &value
			(*output.Reservation)["modification_status"] = value
		case "reservation_net_amortized_upfront_cost_for_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetAmortizedUpfrontCostForUsage = &f
				(*output.Reservation)["net_amortized_upfront_cost_for_usage"] = f
			}
		case "reservation_net_amortized_upfront_fee_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetAmortizedUpfrontFeeForBillingPeriod = &f
				(*output.Reservation)["net_amortized_upfront_fee_for_billing_period"] = f
			}
		case "reservation_net_effective_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetEffectiveCost = &f
				(*output.Reservation)["net_effective_cost"] = f
			}
		case "reservation_net_recurring_fee_for_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetRecurringFeeForUsage = &f
				(*output.Reservation)["net_recurring_fee_for_usage"] = f
			}
		case "reservation_net_unused_amortized_upfront_fee_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetUnusedAmortizedUpfrontFeeForBillingPeriod = &f
				(*output.Reservation)["net_unused_amortized_upfront_fee_for_billing_period"] = f
			}
		case "reservation_net_unused_recurring_fee":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetUnusedRecurringFee = &f
				(*output.Reservation)["net_unused_recurring_fee"] = f
			}
		case "reservation_net_upfront_value":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationNetUpfrontValue = &f
				(*output.Reservation)["net_upfront_value"] = f
			}
		case "reservation_normalized_units_per_reservation":
			output.ReservationNormalizedUnitsPerReservation = &value
			(*output.Reservation)["normalized_units_per_reservation"] = value
		case "reservation_number_of_reservations":
			output.ReservationNumberOfReservations = &value
			(*output.Reservation)["number_of_reservations"] = value
		case "reservation_recurring_fee_for_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationRecurringFeeForUsage = &f
				(*output.Reservation)["recurring_fee_for_usage"] = f
			}
		case "reservation_reservation_a_r_n":
			output.ReservationArn = &value
			(*output.Reservation)["reservation_arn"] = value
		case "reservation_start_time":
			if t, err := helpers.ParseTime(value); err == nil {
				output.ReservationStartTime = &t
				(*output.Reservation)["start_time"] = t
			}
		case "reservation_subscription_id":
			output.ReservationSubscriptionId = &value
			(*output.Reservation)["subscription_id"] = value
		case "reservation_total_reserved_normalized_units":
			output.ReservationTotalReservedNormalizedUnits = &value
			(*output.Reservation)["total_reserved_normalized_units"] = value
		case "reservation_total_reserved_units":
			output.ReservationTotalReservedUnits = &value
			(*output.Reservation)["total_reserved_units"] = value
		case "reservation_units_per_reservation":
			output.ReservationUnitsPerReservation = &value
			(*output.Reservation)["units_per_reservation"] = value
		case "reservation_unused_amortized_upfront_fee_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationUnusedAmortizedUpfrontFeeForBillingPeriod = &f
				(*output.Reservation)["unused_amortized_upfront_fee_for_billing_period"] = f
			}
		case "reservation_unused_normalized_unit_quantity":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationUnusedNormalizedUnitQuantity = &f
				(*output.Reservation)["unused_normalized_unit_quantity"] = f
			}
		case "reservation_unused_quantity":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationUnusedQuantity = &f
				(*output.Reservation)["unused_quantity"] = f
			}
		case "reservation_unused_recurring_fee":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.ReservationUnusedRecurringFee = &f
				(*output.Reservation)["unused_recurring_fee"] = f
			}
		case "reservation_upfront_value":
			if i, err := strconv.ParseInt(value, 10, 64); err == nil {
				output.ReservationUpfrontValue = &i
				(*output.Reservation)["upfront_value"] = i
			}
		case "resource_tags":
			var tags map[string]interface{}
			if err := json.Unmarshal([]byte(value), &tags); err == nil {
				output.ResourceTags = &tags
			}

		case "savings_plan_amortized_upfront_commitment_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanAmortizedUpfrontCommitmentForBillingPeriod = &f
			}
		case "savings_plan_end_time":
			if t, err := helpers.ParseTime(value); err == nil {
				output.SavingsPlanEndTime = &t
			}
		case "savings_plan_instance_type_family":
			output.SavingsPlanInstanceTypeFamily = &value
		case "savings_plan_net_amortized_upfront_commitment_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanNetAmortizedUpfrontCommitmentForBillingPeriod = &f
			}
		case "savings_plan_net_recurring_commitment_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanNetRecurringCommitmentForBillingPeriod = &f
			}
		case "savings_plan_net_savings_plan_effective_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanNetSavingsPlanEffectiveCost = &f
			}
		case "savings_plan_offering_type":
			output.SavingsPlanOfferingType = &value
		case "savings_plan_payment_option":
			output.SavingsPlanPaymentOption = &value
		case "savings_plan_purchase_term":
			output.SavingsPlanPurchaseTerm = &value
		case "savings_plan_recurring_commitment_for_billing_period":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanRecurringCommitmentForBillingPeriod = &f
			}
		case "savings_plan_region":
			output.SavingsPlanRegion = &value
		case "savings_plan_savings_plan_a_r_n":
			output.SavingsPlanSavingsPlanARN = &value
		case "savings_plan_savings_plan_effective_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SavingsPlanSavingsPlanEffectiveCost = &f
			}
		case "savings_plan_savings_plan_rate":
			output.SavingsPlanSavingsPlanRate = &value
		case "savings_plan_start_time":
			if t, err := helpers.ParseTime(value); err == nil {
				output.SavingsPlanStartTime = &t
			}
		case "savings_plan_total_commitment_to_date":
			output.SavingsPlanTotalCommitmentToDate = &value
		case "savings_plan_used_commitment":
			output.SavingsPlanUsedCommitment = &value
		case "split_line_item_actual_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemActualUsage = &f
			}
		case "split_line_item_net_split_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemNetSplitCost = &f
			}
		case "split_line_item_net_unused_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemNetUnusedCost = &f
			}
		case "split_line_item_parent_resource_id":
			output.SplitLineItemParentResourceId = &value
		case "split_line_item_public_on_demand_split_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemPublicOnDemandSplitCost = &f
			}
		case "split_line_item_public_on_demand_unused_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemPublicOnDemandUnusedCost = &f
			}
		case "split_line_item_reserved_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemReservedUsage = &f
			}
		case "split_line_item_split_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemSplitCost = &f
			}
		case "split_line_item_split_usage":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemSplitUsage = &f
			}
		case "split_line_item_split_usage_ratio":
			output.SplitLineItemSplitUsageRatio = &value
		case "split_line_item_unused_cost":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.SplitLineItemUnusedCost = &f
			}
		default:
			switch {
			case strings.HasPrefix(field, "product_"):
				(*output.Product)[strings.Replace(field, "product_", "", 1)] = value
			case strings.HasPrefix(field, "cost_category_"):
				(*output.CostCategory)[strings.Replace(field, "cost_category_", "", 1)] = value
			case strings.HasPrefix(field, "reservation_"):
				(*output.Reservation)[strings.Replace(field, "reservation_", "", 1)] = value
			case strings.HasPrefix(field, "resource_tags_"):
				(*output.ResourceTags)[strings.Replace(field, "resource_tags_", "", 1)] = value
			case strings.HasPrefix(field, "discount_"):
				(*output.Discount)[strings.Replace(field, "discount_", "", 1)] = value
			}
		}
	}

	return output, nil
}

// OnHeader implementOnHeader so that when the collector is notified of a header row, we get notified
func (c *CostAndUsageReportMapper) OnHeader(header []string) {
	c.headers = header
}
