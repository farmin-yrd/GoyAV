package port

import (
	"context"
	"errors"
	"io"
)

// ByteRepository defines the interface for operations related to managing the byte data of documents.
// This interface abstracts the underlying storage mechanism, which could be a file system or
// an object storage system like AWS S3, Azure Blob Storage, or MinIO.
type ByteRepository interface {
	// Save stores the byte data of a document, identified by a unique ID, into the storage system.
	// The function takes a context to manage timeouts and cancellation, a reader for the data,
	// the size of the data, and the document's ID.
	Save(ctx context.Context, data io.Reader, size int64, ID string) error

	// Get retrieves the byte data of a document identified by the given ID.
	// It returns an io.ReadCloser to read the document's data and an error, if any occurred.
	Get(ctx context.Context, ID string) (io.ReadCloser, error)

	// Delete removes the byte data associated with the given document ID from the storage system.
	Delete(ctx context.Context, ID string) error

	// Ping checks the availability or health of the storage system. It is used to verify
	// if the storage system is accessible and functioning correctly.
	Ping() error
}

var (
	// ErrSaveBytesFailed is returned when the Save operation fails.
	ErrSaveBytesFailed = errors.New("failed to save the document's bytes data")

	// ErrGetBytesFailed is returned when the Get operation fails.
	ErrGetBytesFailed = errors.New("failed to save the document's bytes data")

	// ErrDeleteBytesFailed is returned when the Delete operation fails.
	ErrDeleteBytesFailed = errors.New("failed to delete the document's bytes data")

	// ErrByteRepositoryUnavailable is returned when the Ping operation fails to reach the byte repository.
	ErrByteRepositoryUnavailable = errors.New("byte repository is unavailable")
)
