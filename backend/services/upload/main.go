package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pranayjoshi/hoister/backend/services/upload/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
var outDirPath = filepath.Join(dir, "output")

type AWSConfig struct {
	Region      string
	Credentials AWSCredentials
}

type AWSCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func main() {
	// Create the "output" directory if it doesn't exist
	os.MkdirAll(outDirPath, 0755)

	fmt.Println("outDirPath: ", outDirPath)
	PROJECT_ID := os.Getenv("PROJECT_ID")
	BUCKET_REGION := os.Getenv("BUCKET_REGION")
	BUCKET_ACCESS_KEY_ID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	BUCKET_SECRET_ACCESS_KEY := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

	// Check if the environment variables are set
	if PROJECT_ID == "" || BUCKET_REGION == "" || BUCKET_ACCESS_KEY_ID == "" || BUCKET_SECRET_ACCESS_KEY == "" {
		fmt.Println("Error: environment variables not set")
		return
	}

	fmt.Println("Executing build...")
	utils.PublishLog("Build Started...")
	config := AWSConfig{
		Region: BUCKET_REGION,
		Credentials: AWSCredentials{
			AccessKeyId:     BUCKET_ACCESS_KEY_ID,
			SecretAccessKey: BUCKET_SECRET_ACCESS_KEY,
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
		utils.PublishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}

	uploader := s3manager.NewUploader(sess)

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && npm install && npm run build", outDirPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		utils.PublishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}
	utils.PublishLog("Starting to upload")
	err = filepath.Walk(outDirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				utils.PublishLog("Error: " + err.Error())
				return err
			}
			utils.PublishLog("uploading " + path)

			defer file.Close()

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("hoister"),
				Key:    aws.String(fmt.Sprintf("__outputs/%s/%s", PROJECT_ID, strings.TrimPrefix(path, outDirPath))),
				Body:   file,
			})

			if err != nil {
				return err
			}
			utils.PublishLog("uploaded " + path)

			fmt.Println("Successfully uploaded", path)
		}

		return nil
	})

	if err != nil {
		utils.PublishLog("Error: " + err.Error())
		fmt.Println("Error", err)
		return
	}
	utils.PublishLog("Done.. ")
	fmt.Println("Done...")
}
