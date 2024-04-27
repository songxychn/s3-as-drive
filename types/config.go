package types

type S3Config struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
}

type DownloadConfig struct {
	Dir string `json:"dir"`
}

type Config struct {
	S3Config       S3Config       `json:"s3Config"`
	DownloadConfig DownloadConfig `json:"downloadConfig"`
}
