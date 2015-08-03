package s3mgmt

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "fmt"
  "os"
  "log"
)

func ListBuckets(region string) {
  config := aws.NewConfig().WithRegion(region)
  svc := s3.New(config)

  var params *s3.ListBucketsInput
  resp, err := svc.ListBuckets(params)

  if err != nil {
  	if awsErr, ok := err.(awserr.Error); ok {
  		// Generic AWS error with Code, Message, and original error (if any)
  		fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
  		if reqErr, ok := err.(awserr.RequestFailure); ok {
  			// A service error occurred
  			fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
  		}
  	} else {
  		// This case should never be hit, the SDK should always return an
  		// error which satisfies the awserr.Error interface.
  		fmt.Println(err.Error())
  	}
  }

  // Pretty-print the response data.
  fmt.Println(awsutil.Prettify(resp))
}

func ListBucketContents(region, bucket string) {
  config := aws.NewConfig().WithRegion(region)
  svc := s3.New(config)

  params := &s3.ListObjectsInput{
  	Bucket:       aws.String(bucket), // Required
  	// Delimiter:    aws.String("Delimiter"),
  	// EncodingType: aws.String("EncodingType"),
  	// Marker:       aws.String("Marker"),
  	// MaxKeys:      aws.Int64(1),
  	// Prefix:       aws.String("Prefix"),
  }
  resp, err := svc.ListObjects(params)

  if err != nil {
  	if awsErr, ok := err.(awserr.Error); ok {
  		// Generic AWS error with Code, Message, and original error (if any)
  		fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
  		if reqErr, ok := err.(awserr.RequestFailure); ok {
  			// A service error occurred
  			fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
  		}
  	} else {
  		// This case should never be hit, the SDK should always return an
  		// error which satisfies the awserr.Error interface.
  		fmt.Println(err.Error())
  	}
  }

  // Pretty-print the response data.
  fmt.Println(awsutil.Prettify(resp))
}

func BuildUploader(region string) *s3manager.Uploader {
  aws.DefaultConfig.Region = aws.String(region)
  uploader := s3manager.NewUploader(nil)
  return uploader
}

func UploadFile(uploader *s3manager.Uploader, bucket string, path string, contentType string, file *os.File) {
  _, err := uploader.Upload(&s3manager.UploadInput{
      Bucket: &bucket,
      Key:    aws.String(path),
      Body:   file,
      ContentType: &contentType,
  })
  if err != nil {
      log.Fatalln("Failed to upload", path, err)
  }
}
