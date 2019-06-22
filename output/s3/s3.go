package s3

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	session *session.Session
}

//InitOutput initializes S3 with session
func (s *S3) InitOutput() error {

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return err
	}

	s.session = session
	return nil
}

func (s *S3) Run() error {

	file, err := os.Open(os.Getenv("FILE_PATH"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3.New(s.session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("AWS_BUCKET")),
		Key:                aws.String(os.Getenv("FILE_PATH")),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
