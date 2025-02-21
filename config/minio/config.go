package minio

type Minio struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"rootUser"`
	SecretAccessKey string `mapstructure:"rootPassword"`
	BucketName      string `mapstructure:"bucketName"`
}
