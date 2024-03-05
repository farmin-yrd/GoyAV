package byterepo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"goyav/internal/core/port"
	"goyav/pkg/helper"
	"io"
)

// MockByteRepository is a mock implementation of the ByteRepository interface.
// It simulates the behavior of a real repository for testing purposes.
type MockByteRepository struct {
	// simulatedStorage simulates a storage system using a map.
	simulatedStorage map[string][]byte
	isOnline         bool
}

// NewMock creates a new instance of MockByteRepository.
func NewMock() *MockByteRepository {
	return &MockByteRepository{
		simulatedStorage: make(map[string][]byte),
		isOnline:         true,
	}
}

var ErrMockByteRepository = errors.New("MockByteRepository")

// Save simulates the saving of document's byte data.
// It returns ErrSaveFailed error with additional context if the operation fails.
func (m *MockByteRepository) Save(ctx context.Context, data io.Reader, size int64, documentID string) error {

	if !helper.IsValidID(documentID) {
		return fmt.Errorf("%w: %w: invalide id: %q", ErrMockByteRepository, port.ErrSaveBytesFailed, documentID)
	}

	if ctx.Err() != nil {
		return fmt.Errorf("%w: %w: %v", ErrMockByteRepository, port.ErrSaveBytesFailed, ctx.Err())
	}

	data = io.LimitReader(data, size)
	b, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("%w: %w: reading data failed: %v", ErrMockByteRepository, port.ErrSaveBytesFailed, err)
	}
	// Simulate successful save operation.
	m.simulatedStorage[documentID] = b
	return nil
}

// Delete simulates the deletion of document's byte data.
// It returns ErrDeleteFailed error with additional context if the operation fails.
func (m *MockByteRepository) Delete(ctx context.Context, documentID string) error {
	if _, exists := m.simulatedStorage[documentID]; !exists {
		return fmt.Errorf("%w: %w : id not found : id=%q", ErrMockByteRepository, port.ErrDeleteBytesFailed, documentID)
	}

	if ctx.Err() != nil {
		return fmt.Errorf("%w: %w: %v", port.ErrDeleteBytesFailed, port.ErrDeleteBytesFailed, ctx.Err())
	}

	// Simulate successful delete operation.
	delete(m.simulatedStorage, documentID)
	return nil
}

func (m *MockByteRepository) Get(ctx context.Context, ID string) (io.Reader, error) {
	b, exists := m.simulatedStorage[ID]
	if !exists {
		return nil, fmt.Errorf("%w: %w : id not found", ErrMockByteRepository, port.ErrGetBytesFailed)
	}
	return bytes.NewBuffer(b), nil
}

// Ping simulates a check on the storage system.
// It returns ErrPingByteRepositoryFailed if the simulated ping fails.
func (m *MockByteRepository) Ping() error {
	// Simulate a condition that would cause the ping operation to fail.
	if !m.isOnline {
		return fmt.Errorf("%w: %w: reposotory appears to be offline", ErrMockByteRepository, port.ErrByteRepositoryUnavailable)
	}

	return nil
}

func (m *MockByteRepository) Toggle() bool {
	m.isOnline = !m.isOnline
	return m.isOnline
}
