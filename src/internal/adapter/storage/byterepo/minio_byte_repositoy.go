package byterepo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"goyav/internal/core/port"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioByteRepository provides a storage backend using Minio.
type MinioByteRepository struct {
	client     *minio.Client
	bucketName string
}

var ErrMinioByteRepository = errors.New("MinioByteRepository")

// NewMinio creates a new instance of MinioByteRepository.
func NewMinio(host string, port uint64, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinioByteRepository, error) {
	endpoint := fmt.Sprintf("%v:%v", host, port)
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%w : failed to create byte repository: %v", ErrMinioByteRepository, err)
	}

	bucketExists, err := cli.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMinioByteRepository, err)
	}

	if !bucketExists {
		if err = cli.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMinioByteRepository, err)
		}
	}

	return &MinioByteRepository{
		client:     cli,
		bucketName: bucketName,
	}, nil
}

// Save saves an object into the Minio bucket
func (m *MinioByteRepository) Save(ctx context.Context, data io.Reader, size int64, ID string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, ID, io.LimitReader(data, size), size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("%w: %w: %v", ErrMinioByteRepository, port.ErrSaveBytesFailed, err)
	}

	return nil
}

// Delete removes an object from the Minio bucket.
func (m MinioByteRepository) Delete(ctx context.Context, ID string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, ID, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return fmt.Errorf("MinioByteRepository: %w: %v", port.ErrDeleteBytesFailed, err)
	}
	return nil
}

// Get returns an object from the Minio bucket
func (m MinioByteRepository) Get(ctx context.Context, ID string) (io.ReadCloser, error) {
	o, err := m.client.GetObject(ctx, m.bucketName, ID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("MinioByteRepository: %w: %v", port.ErrGetBytesFailed, err)
	}
	return o, nil
}

// Ping checks the availability of the Minio service.
func (m MinioByteRepository) Ping() error {
	if _, err := m.client.ListBuckets(context.Background()); err != nil {
		return fmt.Errorf("MinioByteRepository: %w: %v", port.ErrByteRepositoryUnavailable, err)
	}
	return nil
}
