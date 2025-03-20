package cost_optimization_recommendation

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// Based on AWS Cost Optimization Recommendations columns
// Reference: https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cor-columns.html
type CostOptimizationRecommendation struct {
	schema.CommonFields

	AccountID                                *string                 `json:"account_id,omitempty" parquet:"name=account_id"`
	ActionType                               *string                 `json:"action_type,omitempty" parquet:"name=action_type"`
	CurrencyCode                             *string                 `json:"currency_code,omitempty" parquet:"name=currency_code"`
	CurrentResourceDetails                   *string                 `json:"current_resource_details,omitempty" parquet:"name=current_resource_details"`
	CurrentResourceSummary                   *string                 `json:"current_resource_summary,omitempty" parquet:"name=current_resource_summary"`
	CurrentResourceType                      *string                 `json:"current_resource_type,omitempty" parquet:"name=current_resource_type"`
	EstimatedMonthlyCostAfterDiscount        *float64                `json:"estimated_monthly_cost_after_discount,omitempty" parquet:"name=estimated_monthly_cost_after_discount"`
	EstimatedMonthlyCostBeforeDiscount       *float64                `json:"estimated_monthly_cost_before_discount,omitempty" parquet:"name=estimated_monthly_cost_before_discount"`
	EstimatedMonthlySavingsAfterDiscount     *float64                `json:"estimated_monthly_savings_after_discount,omitempty" parquet:"name=estimated_monthly_savings_after_discount"`
	EstimatedMonthlySavingsBeforeDiscount    *float64                `json:"estimated_monthly_savings_before_discount,omitempty" parquet:"name=estimated_monthly_savings_before_discount"`
	EstimatedSavingsPercentageAfterDiscount  *float64                `json:"estimated_savings_percentage_after_discount,omitempty" parquet:"name=estimated_savings_percentage_after_discount"`
	EstimatedSavingsPercentageBeforeDiscount *float64                `json:"estimated_savings_percentage_before_discount,omitempty" parquet:"name=estimated_savings_percentage_before_discount"`
	ImplementationEffort                     *string                 `json:"implementation_effort,omitempty" parquet:"name=implementation_effort"`
	LastRefreshTimestamp                     *time.Time              `json:"last_refresh_timestamp,omitempty" parquet:"name=last_refresh_timestamp"`
	RecommendationID                         *string                 `json:"recommendation_id,omitempty" parquet:"name=recommendation_id"`
	RecommendationLookbackPeriodInDays       *int                    `json:"recommendation_lookback_period_in_days,omitempty" parquet:"name=recommendation_lookback_period_in_days"`
	RecommendationSource                     *string                 `json:"recommendation_source,omitempty" parquet:"name=recommendation_source"`
	RecommendedResourceDetails               *string                 `json:"recommended_resource_details,omitempty" parquet:"name=recommended_resource_details"`
	RecommendedResourceSummary               *string                 `json:"recommended_resource_summary,omitempty" parquet:"name=recommended_resource_summary"`
	RecommendedResourceType                  *string                 `json:"recommended_resource_type,omitempty" parquet:"name=recommended_resource_type"`
	Region                                   *string                 `json:"region,omitempty" parquet:"name=region"`
	ResourceARN                              *string                 `json:"resource_arn,omitempty" parquet:"name=resource_arn"`
	RestartNeeded                            *bool                   `json:"restart_needed,omitempty" parquet:"name=restart_needed"`
	RollbackPossible                         *bool                   `json:"rollback_possible,omitempty" parquet:"name=rollback_possible"`
	Tags                                     *map[string]interface{} `json:"tags,omitempty" parquet:"name=tags"`
}

func (c *CostOptimizationRecommendation) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"account_id":                                   "The account that the recommendation is for.",
		"action_type":                                  "The type of action you can take by adopting the recommendation.",
		"currency_code":                                "The currency code used for the recommendation.",
		"current_resource_details":                     "The details for the resource in JSON string format.",
		"current_resource_summary":                     "A description of the current resource.",
		"current_resource_type":                        "The type of resource.",
		"estimated_monthly_cost_after_discount":        "The estimated monthly cost of the current resource after discounts. For Reserved Instances and Savings Plans, it refers to the cost for eligible usage.",
		"estimated_monthly_cost_before_discount":       "The estimated monthly cost of the current resource before discounts. For Reserved Instances and Savings Plans, it refers to the cost for eligible usage.",
		"estimated_monthly_savings_after_discount":     "The estimated monthly savings amount for the recommendation after discounts.",
		"estimated_monthly_savings_before_discount":    "The estimated monthly savings amount for the recommendation before discounts.",
		"estimated_savings_percentage_after_discount":  "The estimated savings percentage after discounts relative to the total cost over the cost calculation lookback period.",
		"estimated_savings_percentage_before_discount": "The estimated savings percentage before discounts relative to the total cost over the cost calculation lookback period.",
		"implementation_effort":                        "The effort required to implement the recommendation.",
		"last_refresh_timestamp":                       "The time when the recommendation was last generated.",
		"recommendation_id":                            "The ID for the recommendation.",
		"recommendation_lookback_period_in_days":       "The lookback period that's used to generate the recommendation.",
		"recommendation_source":                        "The source of the recommendation.",
		"recommended_resource_details":                 "The details about the recommended resource in JSON string format.",
		"recommended_resource_summary":                 "A description of the recommended resource.",
		"recommended_resource_type":                    "The resource type of the recommendation.",
		"region":                                       "The AWS Region of the resource.",
		"resource_arn":                                 "The Amazon Resource Name (ARN) of the resource.",
		"restart_needed":                               "Whether or not implementing the recommendation requires a restart.",
		"rollback_possible":                            "Whether or not implementing the recommendation can be rolled back.",
		"tags":                                         "A list of tags associated with the resource for which the recommendation exists.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "The list of resource ARNs associated with a cost and usage recommendation.",
		"tp_index":     "The AWS account ID associated with the recommendation.",
		"tp_timestamp": "The timestamp when the recommendation was last refreshed, in ISO 8601 format.",
	}
}
