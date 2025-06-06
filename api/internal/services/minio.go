package services

import (
	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	MinioClient *minio.Client
}

func NewMinioService(minioClient *minio.Client) *MinioService {
	return &MinioService{
		MinioClient: minioClient,
	}
}

func (ms *MinioService) StoreImage() {

}
