package tables

import (
    "github.com/hashicorp/hcl/v2"
)

type AlbAccessLogTableConfig struct {
    // Required to allow partial decoding
    Remain hcl.Body `hcl:",remain" json:"-"`

    // Optional: Custom log format if not using default
    LogFormat *string `json:"log_format,omitempty" hcl:"log_format,optional"`
    
    // Optional: File pattern for S3 bucket source
    FilePattern *string `json:"file_pattern,omitempty" hcl:"file_pattern,optional"`
    
    // Optional: Timezone for parsing timestamps
    Timezone *string `json:"timezone,omitempty" hcl:"timezone,optional"`
}

func (c *AlbAccessLogTableConfig) Validate() error {
    // By default, use built-in ALB log format
    // Add any additional validation as needed
    return nil
}