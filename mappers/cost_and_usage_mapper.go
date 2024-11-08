package mappers

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// CostAndUsageLogMapper
type CostAndUsageLogMapper struct {
}

// NewCostAndUsageMapper creates a new CostAndUsageLogMapper
func NewCostAndUsageMapper() table.Mapper[rows.CostAndUsageLog] {
	return &CostAndUsageLogMapper{}
}

func (c *CostAndUsageLogMapper) Identifier() string {
	return "cost_and_usage_mapper"
}

func (c *CostAndUsageLogMapper) Map(_ context.Context, a any) ([]rows.CostAndUsageLog, error) {
	slog.Debug(">> Inside Map")
	csvData, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("input data is not of type string")
	}
	r := csv.NewReader(strings.NewReader(csvData))
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return nil, err
	}
	slog.Debug("<< records", "...", records)
	var logs []rows.CostAndUsageLog
	for _, row := range records {
		// TODO - find a better way of doing this parsing
		slog.Debug("<< row", "...", row)
		log := rows.CostAndUsageLog{
			InvoiceID:              &row[0],
			PayerAccountId:         &row[1],
			LinkedAccountId:        &row[2],
			RecordType:             &row[3],
			RecordID:               &row[4],
			BillingPeriodStartDate: parseTime(row[5]),
			BillingPeriodEndDate:   parseTime(row[6]),
			InvoiceDate:            parseTime(row[7]),
			PayerAccountName:       &row[8],
			LinkedAccountName:      &row[9],
			TaxationAddress:        &row[10],
			PayerPONumber:          &row[11],
			ProductCode:            &row[12],
			ProductName:            &row[13],
			SellerOfRecord:         &row[14],
			UsageType:              &row[15],
			Operation:              &row[16],
			RateId:                 &row[17],
			ItemDescription:        &row[18],
			UsageStartDate:         parseTime(row[19]),
			UsageEndDate:           parseTime(row[20]),
			UsageQuantity:          parseFloat(row[21]),
			BlendedRate:            &row[22],
			CurrencyCode:           &row[23],
			CostBeforeTax:          parseFloat(row[24]),
			Credits:                parseFloat(row[25]),
			TaxAmount:              parseFloat(row[26]),
			TaxType:                &row[27],
			TotalCost:              parseFloat(row[28]),
		}
		logs = append(logs, log)
	}
	slog.Debug("<< logs", "...", logs)
	return logs, nil
}

func parseTime(value string) *time.Time {
	if value == "" {
		return nil
	}
	layout := "2006/01/02 15:04:05"
	t, err := time.Parse(layout, value)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return nil
	}
	return &t
}

func parseFloat(value string) *float64 {
	if value == "" {
		return nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println("Error parsing float:", err)
		return nil
	}
	return &f
}
