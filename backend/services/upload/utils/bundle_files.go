package utils

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func BundleFiles(outDirPath string, uploader *s3manager.Uploader, PROJECT_ID string) {

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	filepath.Walk(outDirPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(file, outDirPath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("hoister"),
		Key:    aws.String(fmt.Sprintf("__outputs/%s/output.tar", PROJECT_ID)),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		log.Fatalln("Failed to upload", err)
	}
	PublishLog("Successfully uploaded archive")

	fmt.Println("Successfully uploaded archive")
}
