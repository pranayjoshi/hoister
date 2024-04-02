package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// var ctx = context.Background()

var outDirPath = filepath.Join(filepath.Dir(os.Args[0]), "./outputs")

type AWSConfig struct {
	Region      string
	Credentials AWSCredentials
}

type AWSCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func main() {
	fmt.Println("Executing build...")
	publishLog("Build Started...")
	config := AWSConfig{
		Region: "your-region",
		Credentials: AWSCredentials{
			AccessKeyId:     "your-access-key-id",
			SecretAccessKey: "your-secret-access-key",
		},
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.Credentials.AccessKeyId,
			config.Credentials.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		publishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}

	uploader := s3manager.NewUploader(sess)

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && npm install && npm run build", outDirPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		publishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}
	publishLog("Starting to upload")
	err = filepath.Walk(outDirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		PROJECT_ID := os.Getenv("PROJECT_ID")

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				publishLog("Error: " + err.Error())
				return err
			}
			publishLog("uploading " + path)

			defer file.Close()

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("hoister-outputs"),
				Key:    aws.String(fmt.Sprintf("__outputs/%s/%s", PROJECT_ID, strings.TrimPrefix(path, outDirPath))),
				Body:   file,
			})

			if err != nil {
				return err
			}
			publishLog("uploaded " + path)

			fmt.Println("Successfully uploaded", path)
		}

		return nil
	})

	if err != nil {
		publishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}
	publishLog("Done.. ")
	fmt.Println("Done...")
}
