package utils

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

func getS3LogList(singerKey string) []*string {
	var objectList []*string

	prefix := s3Path()
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(""),
		Prefix: prefix,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, object := range output.Contents {
		if *object.Key == singerKey {
			objectList = append(objectList, object.Key)
		}
	}
	return objectList
}

func formatOjbect(strlog string) (*Alb, error) {
	strings.Replace(strlog, " ", "", -1)
	strings.Replace(strlog, "\n", "", -1)

	re := regexp.MustCompile(`([^ ]*) ([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*)[:-]([0-9]*) ([-.0-9]*) ([-.0-9]*) ([-.0-9]*) (|[-0-9]*) (-|[-0-9]*) ([-0-9]*) ([-0-9]*) \"([^ ]*) (.*) (- |[^ ]*)\" \"([^\"]*)\" ([A-Z0-9-_]+) ([A-Za-z0-9.-]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^\"]*)\" ([-.0-9]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^ ]*)\" \"([^\s]+?)\" \"([^\s]+)\" \"([^ ]*)\" \"([^ ]*)\"`)

	s := re.FindAllStringSubmatch(strlog, -1)

	client_port, err := strconv.Atoi(s[0][5])
	if err != nil {
		return &Alb{}, err
	}

	target_port, err := strconv.Atoi(s[0][7])
	if err != nil {
		return &Alb{}, err
	}
	request_processing_time, err := strconv.ParseFloat(s[0][8], 64)
	if err != nil {
		return &Alb{}, err
	}
	target_processing_time, err := strconv.ParseFloat(s[0][9], 64)
	if err != nil {
		return &Alb{}, err
	}
	response_processing_time, err := strconv.ParseFloat(s[0][10], 64)
	if err != nil {
		return &Alb{}, err
	}
	elb_status_code, err := strconv.Atoi(s[0][11])
	if err != nil {
		return &Alb{}, err
	}
	received_bytes, err := strconv.Atoi(s[0][13])
	if err != nil {
		return &Alb{}, err
	}
	sent_bytes, err := strconv.Atoi(s[0][14])
	if err != nil {
		return &Alb{}, err
	}

	logMap := &Alb{
		s[0][1],
		s[0][2],
		s[0][3],
		s[0][4],
		client_port,
		s[0][6],
		target_port,
		request_processing_time,
		target_processing_time,
		response_processing_time,
		elb_status_code,
		s[0][12],
		received_bytes,
		sent_bytes,
		s[0][15],
		s[0][16],
		s[0][17],
		s[0][18],
		s[0][19],
		s[0][20],
		s[0][21],
		s[0][22],
		s[0][23],
		s[0][24],
		s[0][25],
		s[0][26],
		s[0][27],
		s[0][28],
		s[0][29],
		s[0][30],
		s[0][31],
		s[0][32],
		s[0][33],
	}

	return logMap, nil

}

func GetS3Object(singerKey string) (snippetLog []*Alb) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(""),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "",
				SecretAccessKey: "",
			},
		}))

	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(cfg)

	gzipString := getS3LogList(singerKey)

	for _, key := range gzipString {
		output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(""),
			Key:    key})

		if err != nil {
			fmt.Println(err)
		}

		snippets := unzipOutput(output)

		for _, v := range snippets {

			if len(v) > 10 {
				logMap, err := formatOjbect(v)
				if err != nil {
					continue
				}

				snippetLog = append(snippetLog, logMap)
			}
		}
	}

	return
}

func unzipOutput(output *s3.GetObjectOutput) []string {

	zr, err := gzip.NewReader(output.Body)

	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(zr)

	if err != nil {
		log.Fatal(err)
	}

	snippets := strings.Split(string(b), "\n")

	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}

	return snippets
}
