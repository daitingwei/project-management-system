package min

import (
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strconv"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

var minioClient *MinioClient
var bucketName string
var endpoint string

func InitMinio(cfg *MinioConfig) error {
	client, err := New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.UseSSL)
	if err != nil {
		return err
	}
	minioClient = client
	bucketName = cfg.Bucket
	endpoint = cfg.Endpoint

	ctx := context.Background()
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
	minioClient.c.SetBucketPolicy(ctx, bucketName, policy)
	return nil
}

func GetBucket() string {
	return bucketName
}

func GetMinioClient() (*MinioClient, error) {
	return minioClient, nil
}

func GetEndpoint() string {
	return endpoint
}

type MinioClient struct {
	c *minio.Client
}

func (c *MinioClient) Put(
	ctx context.Context,
	bucketName string,
	fileName string,
	data []byte,
	size int64,
	contentType string,
) (minio.UploadInfo, error) {
	object, err := c.c.PutObject(
		ctx,
		bucketName,
		fileName,
		bytes.NewBuffer(data),
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return object, err
}

func (c *MinioClient) Compose(
	ctx context.Context,
	bucketName string,
	fileName string,
	totalChunks int,
) (minio.UploadInfo, error) {
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: fileName,
	}
	var srcs []minio.CopySrcOptions
	for i := 1; i <= totalChunks; i++ {
		formatInt := strconv.FormatInt(int64(i), 10)
		src := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: fileName + "_" + formatInt,
		}
		srcs = append(srcs, src)
	}
	object, err := c.c.ComposeObject(
		ctx,
		dst,
		srcs...,
	)
	return object, err
}

func (c *MinioClient) Delete(ctx context.Context, bucketName, objectName string) error {
	return c.c.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (c *MinioClient) DeleteObjects(ctx context.Context, bucketName string, objectNames []string) error {
	for _, objectName := range objectNames {
		err := c.c.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MinioClient) Get(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return c.c.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (c *MinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	return c.c.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
}

func New(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{c: minioClient}, err
}
