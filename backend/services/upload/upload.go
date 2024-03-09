package upload

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}))

	uploader := s3manager.NewUploader(sess)

	outDirPath := "./output"

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && npm install && npm run build", outDirPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	err = filepath.Walk(outDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("vercel-clone-outputs"),
				Key:    aws.String(fmt.Sprintf("__outputs/%s/%s", "PROJECT_ID", strings.TrimPrefix(path, outDirPath))),
				Body:   file,
			})

			if err != nil {
				return err
			}

			fmt.Println("Successfully uploaded", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Done...")
}
