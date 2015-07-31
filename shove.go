// go get -u github.com/aws/aws-sdk-go/...
package main

import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/awserr"
import "github.com/aws/aws-sdk-go/aws/awsutil"
import "github.com/aws/aws-sdk-go/service/s3"
import "fmt"

func main() {
  aws.DefaultConfig.Region = aws.String("us-west-2")

  svc := s3.New(nil)

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
