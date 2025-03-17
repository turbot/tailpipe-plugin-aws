package tables

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Caller identity data includes:
// Account ID
// Arn
// UserId
func GetCallerIdentityData() (*sts.GetCallerIdentityOutput, error) {
	callerIdentityData, err := sts.New(session.Must(session.NewSession())).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return callerIdentityData, nil
}

// AwsAkasFromArn will extract key identifiers from an AWS ARN string. For example:
// * the full arn
// * the account ID
// * EC2 instance ID
// * S3 bucket name
// * EC2 volume ID

func AwsAkasFromArn(arn string) []string {
	// Split the ARN into its components.
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return []string{}
	}

	// Extract the service name and the resource descriptor.
	service := parts[2]
	resourceDescriptor := parts[5]
	accountID := parts[4]

	// Initialize a slice to hold the key elements.
	keyElements := []string{arn}
	if accountID != "" {
		keyElements = append(keyElements, accountID)
	}

	// Handle different services.
	switch service {
	case "s3":
		// For S3, the resource descriptor is the bucket name.
		keyElements = append(keyElements, resourceDescriptor)
	case "ec2":
		// For EC2, we need to further parse the resource descriptor.
		if strings.HasPrefix(resourceDescriptor, "instance/") {
			// Extract the instance ID for EC2 instances.
			instanceID := strings.TrimPrefix(resourceDescriptor, "instance/")
			keyElements = append(keyElements, instanceID)
		} else if strings.HasPrefix(resourceDescriptor, "volume/") {
			// Extract the volume ID for EC2 volumes.
			volumeID := strings.TrimPrefix(resourceDescriptor, "volume/")
			keyElements = append(keyElements, volumeID)
		}
	}

	return keyElements
}

func NilIfDash(field *string) *string {
	if field != nil && *field == "-" {
		return nil
	}
	return field
}
// StringToFloat converts a string to a float64 and handles errors
func StringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to float: %v", err)
	}
	return f, nil
}

// StringToTimestamp converts a string to time.Time
func StringToTimestamp(dateStr string, format string) (time.Time, error) {
	parsedTime, err := time.Parse(format, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp from string '%s' using format '%s': %w", dateStr, format, err)
	}
	return parsedTime, nil
}

// StringToMap safely converts a JSON string to map[string]string
func StringToMap(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to map: %w", err)
	}
	return result, nil
}
