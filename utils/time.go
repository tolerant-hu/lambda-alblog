package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

func s3Path() *string {
	spart := getTimeNow()

	s3Path := filepath.Join("", spart)
	fixSlash := s3Path + "/"

	fmt.Printf("it's reading from date: %s \n", fixSlash)
	return &fixSlash
}

func getTimeNow() string {
	now := time.Now()
	dateStr := now.Format("2006/01/02")

	return dateStr
}
