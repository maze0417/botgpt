package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"sync"
)

var s3Client *s3.S3
var once2 sync.Once

type S3 struct {
}

func NewS3() *S3 {
	return &S3{}
}

func getS3Client() *s3.S3 {
	once2.Do(func() {
		// Create a new session with the default AWS configuration
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}

		// Create a Polly client
		s3Client = s3.New(sess)
	})

	return s3Client
}

func (s *S3) Upload(localAudioFilePath string, data []byte) (string, error) {

	s3Client := getS3Client()
	bucketName := "aibaby"
	objectKey := filepath.Base(localAudioFilePath)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"), // 設置為公開可讀取
	})

	if err != nil {
		log.Println("無法上傳 MP3 檔案到 S3:", err)
		return "", err
	}

	// 使用 S3 上傳後的 URL 發送語音訊息
	audioFileURL := "https://" + bucketName + ".s3.amazonaws.com/" + objectKey

	return audioFileURL, nil
}
