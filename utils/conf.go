package utils

type Config struct {
	AwsAccessKeyID     string
	AwsRegion          string
	AwsSecretAccessKey string
	ElkAddresses       string
	ElkIndex           string
	ElkPassword        string
	ElkUsername        string
	S3Bucket           string
	KeyPrefix          *string
	KeyName            string
}
