package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type SecurityHubFinding struct {
	enrichment.CommonFields

	// Top level fields
	Version        *string                 `json:"version,omitempty"`
	ID             *string                 `json:"id,omitempty"`
	DetailType     *string                 `json:"detail_type,omitempty"`
	Source         *string                 `json:"source,omitempty"`
	Account        *string                 `json:"account,omitempty"`
	Time           *time.Time              `json:"time,omitempty"`
	Region         *string                 `json:"region,omitempty"`
	Resources      []*string               `json:"resources,omitempty"`
	DetailFindings *map[string]interface{} `json:"detail" parquet:"name=detail, type=JSON"`
}
