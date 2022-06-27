package main

import (
	"context"
	"fmt"
	"path/filepath"

	"lambda-alblog/utils"

	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type Object struct {
	Key string `json:"key"`
}

type S3 struct {
	Object Object `json:"object"`
}

type Mp struct {
	S3 S3 `json:"s3"`
}

type MyEvent struct {
	Records []Mp `json:"Records"`
}

var config *utils.Config

func init() {
	config = &utils.Config{
		AwsAccessKeyID:     os.Getenv("Aws_AccessKeyID"),
		AwsRegion:          os.Getenv("Aws_Region"),
		AwsSecretAccessKey: os.Getenv("Aws_SecretAccessKey"),
		ElkAddresses:       os.Getenv("ELK_Addresses"),
		ElkIndex:           os.Getenv("ELK_Index"),
		ElkPassword:        os.Getenv("ELK_Password"),
		ElkUsername:        os.Getenv("ELK_Username"),
		S3Bucket:           os.Getenv("S3_Bucket"),
	}
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	key := name.Records[0].S3.Object.Key

	config.KeyName = key
	keyPrefix, _ := filepath.Split(key)
	config.KeyPrefix = &keyPrefix

	snipLog := utils.GetS3Object(config)
	utils.SendMessage(config, snipLog)

	return fmt.Sprintf("%s", name.Records[0].S3.Object.Key), nil
}

func main() {
	lambda.Start(HandleRequest)
}
