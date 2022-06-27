package utils

import (
	"compress/gzip"
	"context"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storager struct {
	C      *Config
	Client *s3.Client
	Key    *string
}

func (s *Storager) ListObject() error {
	output, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.C.S3Bucket),
		Prefix: s.C.KeyPrefix,
	})

	if err != nil {
		return err
	}

	for _, object := range output.Contents {
		if *object.Key == s.C.KeyName {
			s.Key = object.Key
		}
	}

	return nil
}

func (s *Storager) GetObject() ([]*Alb, error) {
	var snippetLog []*Alb

	output, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.C.S3Bucket),
		Key:    s.Key,
	})

	if err != nil {
		return nil, err
	}

	snippets, err := s.unzipGz(output)

	if err != nil {
		return nil, err
	}

	for _, v := range snippets {

		if len(v) > 10 {
			logMap, err := s.formatAlb(v)
			if err != nil {
				continue
			}

			snippetLog = append(snippetLog, logMap)
		}
	}

	return snippetLog, nil
}

func (s *Storager) unzipGz(output *s3.GetObjectOutput) ([]string, error) {

	zr, err := gzip.NewReader(output.Body)

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(zr)

	if err != nil {
		return nil, err
	}

	snippets := strings.Split(string(b), "\n")

	err = zr.Close()
	if err != nil {
		return nil, err
	}

	return snippets, nil
}

func (s *Storager) formatAlb(strlog string) (*Alb, error) {
	strings.Replace(strlog, " ", "", -1)
	strings.Replace(strlog, "\n", "", -1)

	re := regexp.MustCompile(`([^ ]*) ([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*)[:-]([0-9]*) ([-.0-9]*) ([-.0-9]*) ([-.0-9]*) (|[-0-9]*) (-|[-0-9]*) ([-0-9]*) ([-0-9]*) \"([^ ]*) (.*) (- |[^ ]*)\" \"([^\"]*)\" ([A-Z0-9-_]+) ([A-Za-z0-9.-]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^\"]*)\" ([-.0-9]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^ ]*)\" \"([^\s]+?)\" \"([^\s]+)\" \"([^ ]*)\" \"([^ ]*)\"`)

	ss := re.FindAllStringSubmatch(strlog, -1)

	client_port, err := strconv.Atoi(ss[0][5])
	if err != nil {
		return &Alb{}, err
	}

	target_port, err := strconv.Atoi(ss[0][7])
	if err != nil {
		return &Alb{}, err
	}
	request_processing_time, err := strconv.ParseFloat(ss[0][8], 64)
	if err != nil {
		return &Alb{}, err
	}
	target_processing_time, err := strconv.ParseFloat(ss[0][9], 64)
	if err != nil {
		return &Alb{}, err
	}
	response_processing_time, err := strconv.ParseFloat(ss[0][10], 64)
	if err != nil {
		return &Alb{}, err
	}
	elb_status_code, err := strconv.Atoi(ss[0][11])
	if err != nil {
		return &Alb{}, err
	}
	received_bytes, err := strconv.Atoi(ss[0][13])
	if err != nil {
		return &Alb{}, err
	}
	sent_bytes, err := strconv.Atoi(ss[0][14])
	if err != nil {
		return &Alb{}, err
	}

	logMap := &Alb{
		ss[0][1],
		ss[0][2],
		ss[0][3],
		ss[0][4],
		client_port,
		ss[0][6],
		target_port,
		request_processing_time,
		target_processing_time,
		response_processing_time,
		elb_status_code,
		ss[0][12],
		received_bytes,
		sent_bytes,
		ss[0][15],
		ss[0][16],
		ss[0][17],
		ss[0][18],
		ss[0][19],
		ss[0][20],
		ss[0][21],
		ss[0][22],
		ss[0][23],
		ss[0][24],
		ss[0][25],
		ss[0][26],
		ss[0][27],
		ss[0][28],
		ss[0][29],
		ss[0][30],
		ss[0][31],
		ss[0][32],
		ss[0][33],
	}

	return logMap, nil

}
