package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lambda-alblog/utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-elasticsearch/v7"
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

var client *s3.Client
var es *elasticsearch.Client
var con *utils.Config

func init() {
	con = &utils.Config{
		AwsAccessKeyID:     os.Getenv("Aws_AccessKeyID"),
		AwsRegion:          os.Getenv("Aws_Region"),
		AwsSecretAccessKey: os.Getenv("Aws_SecretAccessKey"),
		ElkAddresses:       os.Getenv("ELK_Addresses"),
		ElkIndex:           os.Getenv("ELK_Index"),
		ElkPassword:        os.Getenv("ELK_Password"),
		ElkUsername:        os.Getenv("ELK_Username"),
		S3Bucket:           os.Getenv("S3_Bucket"),
	}

	cfgs, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(con.AwsRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     con.AwsAccessKeyID,
				SecretAccessKey: con.AwsSecretAccessKey,
			},
		}))

	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(cfgs)

	cfge := elasticsearch.Config{
		Addresses: []string{con.ElkAddresses},
		Username:  con.ElkUsername,
		Password:  con.ElkPassword,
	}

	es, err = elasticsearch.NewClient(cfge)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	key := name.Records[0].S3.Object.Key

	con.KeyName = key
	keyPrefix, _ := filepath.Split(key)
	con.KeyPrefix = &keyPrefix

	s3In := &utils.Storager{
		C:      con,
		Client: client,
	}

	err := s3In.ListObject()
	if err != nil {
		fmt.Println(err)
	}

	snippetLog, err := s3In.GetObject()
	if err != nil {
		fmt.Println(err)
	}

	esIn := &utils.Elasticer{
		C:  con,
		Es: es,
	}

	err = esIn.SyncIndex(snippetLog)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s", name.Records[0].S3.Object.Key), nil
}

func main() {
	lambda.Start(HandleRequest)
}
