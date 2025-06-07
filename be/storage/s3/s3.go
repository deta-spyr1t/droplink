package s3

import (
	"bytes"
	"droplink-be/storage"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	Client *s3.S3
	Bucket string
}

func NewS3Storage(region, bucket string) (storage.Storage, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}
	return &S3Storage{
		Client: s3.New(sess),
		Bucket: bucket,
	}, nil
}

func (s *S3Storage) Save(file multipart.File, filename string) (string, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	return filename, err
}

func (s *S3Storage) Load(path string) ([]byte, string, error) {
	out, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, "", err
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	return data, path, err
}

func (s *S3Storage) Delete(path string) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	return err
}
