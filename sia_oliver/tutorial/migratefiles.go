package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Function to upload a file to the renterd API
func uploadFileToRenterd(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new multipart buffer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new form file field
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return err
	}

	// Read the file content and write it to the form field
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create a new HTTP request to the renterd API endpoint
	request, err := http.NewRequest(http.MethodPut, "https://renterd-mock.sia.tools/api/worker/objects/"+filePath, body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.SetBasicAuth("", "<auth-token>") // Replace with the actual authentication token

	// Make the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("file upload failed with status code: %d", response.StatusCode)
	}

	return nil
}

// Function to list files in an AWS bucket
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

	// Iterate over the list of files and upload each file to renterd
	for _, file := range files {
		err := uploadFileToRenterd(file)
		if err != nil {
			log.Println("Failed to upload file", file, ":", err)
		} else {
			log.Println("File", file, "uploaded successfully")
		}
	}
}
