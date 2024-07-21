package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func BundleFiles(outDirPath string, uploader *s3manager.Uploader, PROJECT_ID string) {

	err := filepath.Walk(outDirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				PublishLog(PROJECT_ID, "Error: "+err.Error())
				return err
			}
			PublishLog(PROJECT_ID, "uploading "+path)

			defer file.Close()

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("hoister"),
				Key:    aws.String(fmt.Sprintf("__outputs/%s/%s", PROJECT_ID, strings.TrimPrefix(path, outDirPath))),
				Body:   file,
			})

			if err != nil {
				return err
			}
			PublishLog(PROJECT_ID, "uploaded "+path)

			fmt.Println(PROJECT_ID, "Successfully uploaded", path)
		}

		return nil
	})
	if err != nil {
		PublishLog(PROJECT_ID, "Error: "+err.Error())
		fmt.Println("Error", err)
		return
	}
	PublishLog(PROJECT_ID, "Successfully uploaded archive")

	fmt.Println("Successfully uploaded archive")
}
