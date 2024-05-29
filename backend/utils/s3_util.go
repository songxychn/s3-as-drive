package utils

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3Client *minio.Client = nil

func init() {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}
	s3Client, err = minio.New(config.S3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3Config.AccessKey, config.S3Config.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		panic(err)
	}
}
func GetS3Client() *minio.Client {
	return s3Client
}
