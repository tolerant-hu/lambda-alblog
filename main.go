package main

import (
	"context"
	"fmt"

	"lambda-alblog/utils"

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

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	key := name.Records[0].S3.Object.Key

	snipLog := utils.GetS3Object(key)
	utils.SendMessage(snipLog)

	return fmt.Sprintf("Hello %s!", name.Records[0].S3.Object.Key), nil
}

func main() {
	lambda.Start(HandleRequest)
}
