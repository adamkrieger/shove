// go get -u github.com/aws/aws-sdk-go/...
package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/codegangsta/cli"
  "os"
  "fmt"
)

func list(region string) {
  aws.DefaultConfig.Region = aws.String(region)

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

func main() {
  app := cli.NewApp()
  app.Name = "shove"
  app.Usage = "Manage and push files to an S3 bucket."
  app.Commands = []cli.Command{
    {
      Name: "list",
      Aliases: []string{"l"},
      Usage: "List available buckets.",
      Flags: []cli.Flag {
        cli.StringFlag{
          Name: "region, r",
          Value: "us-west-2",
          Usage: "Region to communicate with.",
        },
      },
      Action: func(c *cli.Context) {
        list(c.String("region"))
      },
    },
  }

  app.Run(os.Args)
}
