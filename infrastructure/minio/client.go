package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewClient(host, accessKey, secretKey, bucketName string, useSSL bool) (*Client, error) {
	fmt.Println("Connecting to Minio on", host)
	fmt.Println("Using access key:", accessKey)
	fmt.Println("Using access key:", secretKey)
	endpoint := host
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %v", err)
	}

	found, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	fmt.Printf("Bucket exists: %v\n", found)

	if !found {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket '%s': %v", bucketName, err)
		}
	}

	return &Client{
		MinioClient: client,
		BucketName:  bucketName,
	}, nil
}
