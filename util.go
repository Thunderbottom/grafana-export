package main

import (
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mholt/archiver/v3"
)

// compress is a function that creates a gzipped
// archive for the specified filepath
func compress(filepath string) (string, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// the filepath specified does not exist
		return "", err
	}

	now := time.Now()
	// create a timestamped filename.
	// dashboards/ => dashboards-20060102150405.tar.gz
	arcFn := strings.TrimSuffix(filepath, "/") + "-" + now.Format("20060102150405") + ".tar.gz"

	tarGZ := archiver.NewTarGz()
	if err := tarGZ.Archive([]string{filepath}, arcFn); err != nil {
		return "", err
	}

	return arcFn, nil
}

// backup is a function that takes a file and uploads
// it to the specified s3 bucket
func backup(filepath, bucket, key string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("private"),
		Body:   f,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
