package tables

import (
	"strings"
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