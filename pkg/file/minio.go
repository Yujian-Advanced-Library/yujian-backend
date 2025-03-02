package file

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"yujian-backend/pkg/log"
)

// CreateBucket 创建存储桶
func (client *MinioClient) CreateBucket(ctx context.Context, bucketName string) error {
	// 检查存储桶是否存在
	exists, err := client.inner.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		// 创建存储桶
		err = client.inner.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadFile 上传文件到 MinIO
func (client *MinioClient) UploadFile(ctx context.Context, bucketName, objectName, filePath string) error {
	// 上传文件
	info, err := client.inner.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

// DownloadFile 从 MinIO 下载文件
func (client *MinioClient) DownloadFile(ctx context.Context, bucketName, objectName, filePath string) (string, error) {
	// 下载文件
	err := client.inner.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	fmt.Printf("Successfully downloaded %s to %s\n", objectName, filePath)
	return filePath, nil
}

func (client *MinioClient) FetchFile(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	object, err := client.inner.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.GetLogger().Errorf("Failed to get object: %v", err)
		return nil, err
	}

	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			log.GetLogger().Warnf("Failed to close object: %v", err)
		}
	}(object)

	// 读取, 每次读取4KB
	bufferSize := 4096
	buffer := make([]byte, bufferSize)
	var result []byte
	for {
		n, readErr := object.Read(buffer)
		if n > 0 {
			result = append(result, buffer[:n]...)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			// 其他错误，记录日志并返回
			log.GetLogger().Fatalf("Failed to read file content: %v", readErr)
			return nil, readErr
		}
	}
	return result, nil
}
