package file

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"sync"
	"yujian-backend/pkg/log"
)

var (
	minioClient *MinioClient
	once        sync.Once
)

type MinioClient struct {
	inner *minio.Client
}

func InitMinio() {
	once.Do(func() {
		if client, err := newMinIOClient("http://127.0.0.1:9000", "minioadmin", "minioadmin"); err != nil {
			log.GetLogger().Fatalf("Failed to initialize MinIO client: %v", err)
		} else {
			minioClient = &MinioClient{inner: client}
			log.GetLogger().Info("Successfully initialized MinIO client...")
		}
	})
}

func GetMinioClient() *MinioClient {
	return minioClient
}

func newMinIOClient(endpoint, accessKeyID, secretAccessKey string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // 如果使用 HTTPS 则设置为 true
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
