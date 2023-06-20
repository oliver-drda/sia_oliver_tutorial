package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listFilesInBucket(bucketName string) ([]string, error) {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("<region>"), // Replace with your desired region
	})
	if err != nil {
		return nil, err
	}

	// Create a new S3 service client
	svc := s3.New(sess)

	// Prepare the input parameters for listing objects in the bucket
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	// Perform the list operation
	resp, err := svc.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	// Collect the file names
	var files []string
	for _, obj := range resp.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

func main() {
	// Replace "<your-bucket-name>" with the actual name of your AWS bucket
	files, err := listFilesInBucket("<your-bucket-name>")
	if err != nil {
		log.Fatal(err)
	}

	// Print the list of files
	for _, file := range files {
		fmt.Println(file)
	}
}
