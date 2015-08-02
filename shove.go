package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/codegangsta/cli"
  "os"
  "log"
  "path/filepath"
  "fmt"
)

func listobjects(region, bucket string) {
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

func list(region string) {
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

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if !info.IsDir() {
        f <- path
    }
    return nil
}

func upload(region, bucket, directory string) {
  fmt.Println(bucket)
  prefix := ""

  walker := make(fileWalk)
  go func() {
      // Gather the files to upload by walking the path recursively.
      if err := filepath.Walk(directory, walker.Walk); err != nil {
          log.Fatalln("Walk failed:", err)
      }
      close(walker)
  }()

  // For each file found walking upload it to S3.
  aws.DefaultConfig.Region = aws.String(region)
  uploader := s3manager.NewUploader(nil)
  for path := range walker {
      rel, err := filepath.Rel(directory, path)
      if err != nil {
          log.Fatalln("Unable to get relative path:", path, err)
      }
      file, err := os.Open(path)
      if err != nil {
          log.Println("Failed opening file", path, err)
          continue
      }
      defer file.Close()
      _, err = uploader.Upload(&s3manager.UploadInput{
          Bucket: &bucket,
          Key:    aws.String(filepath.Join(prefix, rel)),
          Body:   file,
      })
      if err != nil {
          log.Fatalln("Failed to upload", path, err)
      }
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "shove"
  app.Usage = "Manage and push files to an S3 bucket."
  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "region, r",
      Value: "us-east-1",
      Usage: "Region to communicate with.",
    },
  }
  app.Commands = []cli.Command{
    {
      Name: "list",
      Aliases: []string{"l"},
      Usage: "List available buckets.",
      Action: func(c *cli.Context) {
        list(c.GlobalString("region"))
      },
    },
    {
      Name: "contents",
      Aliases: []string{"c"},
      Usage: "List the contents of a bucket.",
      Flags: []cli.Flag {
        cli.StringFlag{
          Name: "bucket, b",
          Usage: "Name of the bucket.",
        },
      },
      Action: func(c *cli.Context) {
        listobjects(c.GlobalString("region"), c.String("bucket"))
      },
    },
    {
      Name: "push",
      Aliases: []string{"p"},
      Usage: "Pushes the contents to the bucket.",
      Flags: []cli.Flag {
        cli.StringFlag{
          Name: "directory, d",
          Usage: "Path containing files to be uploaded.",
        },
        cli.StringFlag{
          Name: "bucket, b",
          Usage: "Name of the bucket.",
        },
      },
      Action: func(c *cli.Context) {
        // ./shove push -d "./test" -b "yourbucketname"
        bucket := c.String("bucket")
        if(bucket != ""){
          upload(c.GlobalString("region"), bucket, c.String("directory"))
        }else{
          fmt.Println("Bucket cannot be blank.")
        }
      },
    },
  }

  app.Run(os.Args)
}
