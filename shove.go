package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/codegangsta/cli"
  "github.com/fatih/color"
  "os"
  "log"
  "path/filepath"
  "fmt"
  "./s3mgmt"
)

func getContentType(extension string) string {
  switch extension {
  case ".html":
    return "text/html"
  case ".xml":
    return "text/xml"
  case ".css":
    return "text/css"
  }
  return "binary/octet-stream"
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
  skipColor := color.New(color.FgRed).SprintFunc()

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

  uploader := s3mgmt.BuildUploader(region)

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
      fileInfo, err := file.Stat()
      if err != nil {
        log.Println("Failed getting file stats", path, err)
      }

      fileEmpty := fileInfo.Size() == 0

      defer file.Close()
      contentType := getContentType(filepath.Ext(path))
      destinationPath := filepath.Join(prefix, rel)
      fileName := filepath.Base(path)

      fmt.Println(prefix, rel, contentType, fileInfo.Size(), fileEmpty, fileName, destinationPath)

      if fileEmpty {
        fmt.Printf("%s", skipColor("Skipped (empty)"))
      } else if fileName == ".DS_Store" {
        fmt.Printf("%s", skipColor("Skipped (.DSStore)"))
      } else {
        s3mgmt.UploadFile(uploader, bucket, destinationPath, contentType, file)
      }
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "shove"
  app.Version = "0.1.0"
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
        s3mgmt.ListBuckets(c.GlobalString("region"))
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
        s3mgmt.ListBucketContents(c.GlobalString("region"), c.String("bucket"))
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
