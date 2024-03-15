package adapters

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gwiyeomgo/adapters/config"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
)

var (
	awsS3AdapterOnce     sync.Once
	awsS3AdapterInstance *awsS3Adapter
)

func AwsS3Adapter() *awsS3Adapter {
	awsS3AdapterOnce.Do(func() {
		awsS3AdapterInstance = &awsS3Adapter{}
	})

	return awsS3AdapterInstance
}

type awsS3Adapter struct {
}

func (adapter awsS3Adapter) UploadFile(path string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Config.AwsS3.Region),
		Credentials: credentials.NewStaticCredentials(
			config.Config.AwsS3.AccessKeyId,     // id
			config.Config.AwsS3.SecretAccessKey, // secret
			""),                                 // token can be left blank for now
	})
	if err != nil {
		return "", err
	}

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return "", err
	}

	uploadFileName := fmt.Sprintf("%v/%v", path, uuid.Must(u2, nil).String()+filepath.Ext(fileHeader.Filename))

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(config.Config.AwsS3.Bucket),
		Key:                  aws.String(uploadFileName),
		ACL:                  aws.String("public-read"), // private <- only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	if err != nil {
		return "", err
	}

	accessUrl := fmt.Sprintf("%v/%v", config.Config.AwsS3.HttpEndPoint, uploadFileName)
	return accessUrl, nil
}
