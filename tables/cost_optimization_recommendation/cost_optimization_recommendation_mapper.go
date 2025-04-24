package cost_optimization_recommendation

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/error_types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type CostOptimizationRecommendationMapper struct {
	headers []string
}

func NewCostOptimizationRecommendationMapper() *CostOptimizationRecommendationMapper {
	return &CostOptimizationRecommendationMapper{}
}

func (m *CostOptimizationRecommendationMapper) Identifier() string {
	return "cost_optimization_recommendation_mapper"
}

// OnHeader implements mappers.HeaderHandler so that when the collector is notified of a header, we set headers
func (m *CostOptimizationRecommendationMapper) OnHeader(header []string) {
	m.headers = header
}

func (m *CostOptimizationRecommendationMapper) Map(_ context.Context, a any, opts ...mappers.MapOption[*CostOptimizationRecommendation]) (*CostOptimizationRecommendation, error) {
	var input []byte

	// apply opts
	for _, opt := range opts {
		if opt != nil {
			opt(m)
		}
	}

	// validate input type
	switch v := a.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	default:
		slog.Error("CostOptimizationRecommendationMapper.Map failed to map row due to invalid type", "expected", "[]byte or string", "got", v)
		return nil, error_types.NewRowErrorWithMessage("unable to map row, invalid type received")
	}

	// read CSV line
	reader := csv.NewReader(bytes.NewReader(input))
	record, err := reader.Read()
	if err != nil {
		slog.Error("CostOptimizationRecommendationMapper.Map failed to read CSV line", "error", err)
		return nil, error_types.NewRowErrorWithMessage("failed to read log line")
	}

	// validate header/value count
	if len(record) != len(m.headers) {
		slog.Error("CostOptimizationRecommendationMapper.Map failed to map row due to header/value count mismatch", "expected", len(m.headers), "got", len(record))
		return nil, error_types.NewRowErrorWithMessage("row fields doesn't match count of headers")
	}

	// create a new CostOptimizationRecommendation object, with initialized maps
	output := NewCostOptimizationRecommendation()

	// non-standard time format
	timeFormat := "Mon Jan 2 15:04:05 UTC 2006"

	// map to struct
	for i, value := range record {
		field := strings.ToLower(m.headers[i])

		switch field {
		case "account_id":
			output.AccountID = &value
		case "action_type":
			output.ActionType = &value
		case "currency_code":
			output.CurrencyCode = &value
		case "current_resource_details":
			var currentResourceDetails map[string]interface{}
			if err := json.Unmarshal([]byte(value), &currentResourceDetails); err == nil {
				output.CurrentResourceDetails = &currentResourceDetails
			}
		case "current_resource_summary":
			output.CurrentResourceSummary = &value
		case "current_resource_type":
			output.CurrentResourceType = &value
		case "estimated_monthly_cost_after_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedMonthlyCostAfterDiscount = &f
			}
		case "estimated_monthly_cost_before_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedMonthlyCostBeforeDiscount = &f
			}
		case "estimated_monthly_savings_after_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedMonthlySavingsAfterDiscount = &f
			}
		case "estimated_monthly_savings_before_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedMonthlySavingsBeforeDiscount = &f
			}
		case "estimated_savings_percentage_after_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedSavingsPercentageAfterDiscount = &f
			}
		case "estimated_savings_percentage_before_discount":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				output.EstimatedSavingsPercentageBeforeDiscount = &f
			}
		case "implementation_effort":
			output.ImplementationEffort = &value
		case "last_refresh_timestamp":
			if t, err := time.Parse(timeFormat, value); err == nil {
				output.LastRefreshTimestamp = &t
			}
		case "recommendation_id":
			output.RecommendationID = &value
		case "recommendation_lookback_period_in_days":
			if i, err := strconv.Atoi(value); err == nil {
				output.RecommendationLookbackPeriodInDays = &i
			}
		case "recommendation_source":
			output.RecommendationSource = &value
		case "recommended_resource_details":
			var recommendedResourceDetails map[string]interface{}
			if err := json.Unmarshal([]byte(value), &recommendedResourceDetails); err == nil {
				output.RecommendedResourceDetails = &recommendedResourceDetails
			}
		case "recommended_resource_summary":
			output.RecommendedResourceSummary = &value
		case "recommended_resource_type":
			output.RecommendedResourceType = &value
		case "region":
			output.Region = &value
		case "resource_arn":
			output.ResourceARN = &value
		case "restart_needed":
			if b, err := strconv.ParseBool(value); err == nil {
				output.RestartNeeded = &b
			}
		case "rollback_possible":
			if b, err := strconv.ParseBool(value); err == nil {
				output.RollbackPossible = &b
			}
		case "tags":
			tags := make(map[string]string)
			if err := json.Unmarshal([]byte(value), &tags); err == nil {
				output.Tags = &tags
			}
		}

	}

	return output, nil
}
