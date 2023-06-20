# Migrating Files from an AWS Bucket to `renterd`

<span class="highlight">Go</span>
<span class="highlight">macOS</span>
<span class="highlight">AWSBucket</span>
<span class="highlight">macOS</span>
<span class="highlight">renterd</span>

This tutorial will guide you through the process of building a utility in Go that retrieve files from an AWS bucket to `renterd`. We will cover the following topics:

1. Installing the AWS SDK
2. Setting up AWS credentials
3. Preparing your Bucket
4. Listing files in an AWS bucket
5. Migrating files from the AWS bucket to `renterd`

### Prerequisites 

#### Brew

To install `Homebrew` on macOS, open the Terminal and run following command and press Enter:

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```
Then run this command to check the lastest version has be installed succesfully:

```
brew --version
```

#### AWS CLI

To install `AWS CLI` on macOS, open the Terminal and run following command and press Enter:

```
brew install awscli
```
Then run this command to check the lastest version has be installed succesfully:

```
aws --version
```

```admonish warning title="Note"
For the purpose of this tutorial, we will be demonstrating the process in Go. Go is the native language of the Sia platform, but feel free to adapt the concepts and logic presented here to your preferred programming language of choice by visiting [here](https://docs.aws.amazon.com) to find AWS software development kits (SDKs).
```
#### Go

1. Before we begin, please make sure you have Go installed on your machine. You can download and [install Go](https://golang.org) from the official website. Go is available on the following operating systems and versions:

- __Microsoft Windows__: Windows 7 or later, Intel 64-bit processor.
- __Apple macOS (ARM64)__: macOS 11 or later, Apple 64-bit processor.
- __Apple macOS (x86-64)__: macOS 10.13 or later, intel 64-bit processor.
- __Linux__: 2.6.32 or later, intel 64-bit processor.

2. Install Go and once the download is complete, open the downloaded package file `(.pkg)` and follow the instructions to install Go on your system. The installer will guide you through the installation process.

#### Tutuorial Folder

Create a new directory for your `tutorial`, and `path/to/tutorial` then run:

```
go mod init your-project-name
```

### 1. Installing the AWS SDK

```admonish info title="Here's a helping hand!"
If you're new to AWS (Amazon Workflow Services), these helpful links:

- Familiarise yourself with [AWS](https://aws.amazon.com/what-is-aws/) and the [AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/api/#:~:text=The%20AWS%20SDK%20for%20Go,against%20a%20web%20service%20interface.).
- Basic understanding of [HTTP requests](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_query-requests.html).

```

To interact with AWS services, we need to install the AWS SDK for Go. The SDK provides us with the necessary libraries and tools to work with AWS services, including Amazon S3.

Open up your terminal and run the following command:

```shell
go get github.com/aws/aws-sdk-go@latest
```

This command installs the latest AWS SDK for Go and its dependencies.

### 2. Setting up AWS credentials

To access your AWS bucket, you need to set up your AWS credentials. These credentials consist of an Access Key ID and Secret Access Key associated with an IAM user that has the necessary permissions to access the AWS bucket.

 1. Go to your [AWS Management Console](https://console.aws.amazon.com/console/home) and sign in. If you haven't got an account create a new AWS account by completing the 5-step account creation process.

```admonish warning title="Warning"
Your credit/debit card information will be requested in the account creation process so keep this handy.
```

2. Open the IAM service by searching for "IAM" in the AWS Management Console search bar and selecting the IAM service.

3. Create a new IAM user or use an existing one that has the necessary permissions to access the AWS bucket. 

	a) This can be done firstly by clicking the "Users" link on the dashboard. Once the page, click the "Add Users" button to start creating a new IAM user. Then provide a name for the user in the "User name" field.

	b) When creating a new IAM user or using an existing one, choose an appropriate permission model based on your requirements. For the purpose of this tutorial, choose "Attach existing policies directly", search for the `AmazonS3FullAccess` policy, and check the box next.

	c) Review the settings, click "Next," add any optional tags, proceed to review the user's details, and finally click "Create user."

	d) The `access_key_id` and `secret_access_key` forthe new user will have not be created yet, but you can do this manually. After creating the user, you'll be redirected to the "Users" page. Now click on the newly created user account. Click the "Security credentials" tab and scroll down to the "Access Keys" section and press "Create Access key".

	e) You should now see a page called "Access key best practices &  alternatives". Select the "Command Line Interface" option and remember to check the at the bottom before proceeding by pressing "Next".

	f) You can decide to create tags, but it's not neccessary. Press the "Create Access Key" to generate the `access_key_id` and `secret_access_key`. 


```admonish info title="Make a Note!"
Take note of the `access_key_id` and `secret_access_key` generated for the new user, as these credentials will be needed for programmatically accessing the S3 bucket.
```

4. Return back to your terminal to run the following command to configure the AWS CLI:

```shell
aws configure
```
5. Enter your `access_key_id` and `secret_access_key` when prompted.

6. Choose a default AWS Region (e.g., `us-east-1`) and leave the output format as default.

7. Set the output format as `json`.

```admonish success title=""
Congratulations! The AWS CLI has saved the credentials in a configuration file located in your user's home directory. The AWS CLI configuration is independent of your Go project. It sets up the credentials globally on your machine, and the Go SDK will automatically use them.
```

### 3. Preparing your Bucket.

1. Return back AWS Access Management page of your existing account, on "Users" page, click on the newly created user account. Click the "Security credentials" tab begin the enablement of the User to access the AWS Management Console by pressing the "Enable Console Access" button in the "Console sign-in" section.

2. On the "User" page of the newly created user, there is a link in this format `hhtps://<12-digit-account-id>.signin.aws.amazon.com/console`, that will give you direct access to log into the AWS Access Management specific to the new user you've created.

3. Login using `account_id` (should auto-populate after using the link), `iam_user_name`, and `password`. Once you've logged in, navigate to the S3 service on the "Console Home".

4. Let's start uploading the files in the AWS bucket to see what we have for the listing segment of the tutorial. For the testing purposes, we can create 3 dummy text files (i.e., `file-1.rtf` `file-2.rtf` `file-3.rtf`). Upload these to the bucket for the next steps. 

```admonish info title="Can't see any Buckets?"
If you can't see any buckets in the AWS S3 console, it may be because you haven't created any buckets yet. Here are the steps to create a bucket:

a) Click on the "Create bucket" button.

b) Enter a name for your bucket, choose a region, and click "Create bucket".

c) Your new bucket should now appear in the list of buckets on the S3 console. Once you have created a bucket, you can upload files to it.
```

### 3. Listing files in an AWS bucket

Now we can start by creating a `listfiles.go` file inside our tutorial directory folder and implementing the following by adding the following code:

```go 
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

```

```admonish title="What is this doing?"
In this code snippet, we use the AWS SDK for Go to list objects in the specified AWS bucket. 

- We create an AWS session and configure it with the desired AWS region.
- Using the session, we create an S3 service client.
- We specify the bucket name we want tolist objects from.
- We call the `ListObjectsV2` method on the S3 service client, passing the bucket name in the `Bucket` field of the `ListObjectsV2Input` struct.
- The response contains a list of objects in the bucket, and we iterate over them, printing out the object names.
```

Now run `go run listfiles.go` inside your terminal where the tutorial directory is located.

```admonish success title="Expected Output" 
	mypctest@tests-MacBook-Pro tutorial % go run listfiles.go
	file-1.rtf
	file-2.rtf
	file-3.rtf
```

```admonish success title=""
Congratulations! This part of the tutorial helps us verify that we can successfully connect to the AWS bucket and retrieve a list of objects.
```

### 4. Migrating files from the AWS bucket to `renterd`

Now we can move on to migrating files from the AWS bucket to `renterd`. For the purpose of this tutorial, we will use the mocked `renterd` API endpoints instead of setting up a `renterd` node.

Start by creating a `migratefiles.go` file inside our tutorial directory folder and implementing the following by adding the following code:

```go
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
```

```admonish title="What is this doing?"
This code snippet demonstrates how to migrate a file from the AWS bucket to `renterd` using the mocked `renterd` API endpoints. Here's a breakdown of the code:

- We specify the AWS bucket URL and the `renterd` upload endpoint URL. You should replace these URLs with your actual values.
- We specify the object key (filename) of the file we want to migrate. In this particular case, we made it so it selects them all.
- We use the `http.Get` function to download the object from the AWS bucket.
- We create an upload request to the `renterd` API endpoint by creating an `http.NewRequest` with the `PUT` method and the `renterd` upload URL.
- We set the required HTTP basic authentication header using `req.SetBasicAuth` with an empty username and the provided auth token.
- We send the upload request using `client.Do(req)`.
- Finally, we check the response status. If it's not `http.StatusOK`, we consider the upload as unsuccessful.
```

Now run `go run migratefiles.go` inside your terminal where the tutorial directory is located.

```admonish success title="Expected Output" 
	2023/06/20 18:03:37 File file-1.rtf uploaded successfully
	2023/06/20 18:03:37 File file-2.rtf uploaded successfully
	2023/06/20 18:03:37 File file-3.rtf uploaded successfully
```

```admonish success title="YOU'VE DONE IT🎉" 
Congratulations! 

You've built a utility in Go that allows us to migrate files from an AWS bucket to `renterd`. Here's what you've accomplished, accomplish the following:

**Installation of the AWS SDK**: You've installed the AWS SDK for Go, which provides the necessary libraries and tools to interact with AWS services, including Amazon S3.

**Setting up AWS credentials**: You've configured the AWS credentials, consisting of an Access Key ID and Secret Access Key, to authenticate and access our AWS bucket,ensures that you have the necessary permissions to retrieve files from the bucket.

**Prepared your Bucket**: You've set up an AWS bucket environment for where you can test Listing and Migrate files.

**Listing files in an AWS bucket**: You've utilize the AWS SDK for Go to list the objects/files present in the AWS bucket. This step helps us verify our connectivity to the bucket and retrieve a list of files that can be migrated.

**Migrating files from the AWS bucket to `renterd`**: You've demonstrated the process of migrating a file from the AWS bucket to `renterd`using the appropriate API endpoint. This step showcases the file migration process and confirms that the migration was successful.

```

To get ahold of the complete tutorial, visit [Github](https://github.com/oliver-drda/sia_oliver_tutorial/tree/main/sia_oliver/tutorial)

