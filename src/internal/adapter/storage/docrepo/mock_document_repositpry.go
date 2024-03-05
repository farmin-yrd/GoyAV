package docrepo

import (
	"context"
	"errors"
	"fmt"
	"goyav/internal/core/domain"
	"goyav/internal/core/port"
	"goyav/pkg/helper"
	"time"
)

// MockDocumentRepository is a mock implementation of the DocumentRepository interface.
// It uses an in-memory map to simulate document storage.
type MockDocumentRepository struct {
	documents map[string]*domain.Document
	isOnline  bool
}

var ErrDocumentRepository = errors.New("MockDocumentRepository")

// NewMock creates a new instance of MockDocumentRepository.
func NewMock() *MockDocumentRepository {
	return &MockDocumentRepository{
		documents: make(map[string]*domain.Document),
		isOnline:  true,
	}
}

// Get retrieves a document by its ID.
func (m *MockDocumentRepository) Get(ctx context.Context, id string) (*domain.Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if doc, exists := m.documents[id]; exists {
		return doc, nil
	}
	return nil, fmt.Errorf("%w: %w: id=%q", ErrDocumentRepository, port.ErrDocumentNotFound, id)
}

// Save adds a new document to the repository.
func (m *MockDocumentRepository) Save(ctx context.Context, d *domain.Document) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	doc, _ := m.Get(ctx, d.ID)
	if doc != nil {
		return fmt.Errorf("%w: %w: %w: id=%q", ErrDocumentRepository, port.ErrSaveDocumentFailed, port.ErrDocumentAlreadyExists, doc.ID)
	}
	m.documents[d.ID] = d
	return nil
}

// GetByHash retrieves a document by its hash.
func (m *MockDocumentRepository) GetByHash(ctx context.Context, h string) (*domain.Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if !helper.IsValidHash(h) {
		return nil, fmt.Errorf("%w: invalid hash: %q", ErrDocumentRepository, h)
	}
	for _, doc := range m.documents {
		if doc.Hash == h {
			return doc, nil
		}
	}
	return nil, fmt.Errorf("%w: %w: hash=%q", ErrDocumentRepository, port.ErrDocumentNotFound, h)
}

// Delete removes a document from the repository.
func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := m.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDocumentRepository, port.ErrDeleteDocumentFailed, err)
	}
	delete(m.documents, id)
	return nil
}

// UpdateStatus updates the analysis status and the analysis date of a document.
func (m *MockDocumentRepository) UpdateStatus(ctx context.Context, id string, status domain.AnalysisStatus, analyzedAt time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	doc, err := m.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDocumentRepository, port.ErrUpdateStatusFailed, err)
	}
	doc.Status = status
	doc.AnalyzedAt = analyzedAt
	return nil
}

// Ping checks the availability of the repository.
func (m *MockDocumentRepository) Ping() error {
	// Simulate a condition that would cause the ping operation to fail.
	if !m.isOnline {
		return fmt.Errorf("%w: %w: repository appears to be offline", ErrDocumentRepository, port.ErrDocumentRepositoryUnavailable)
	}
	return nil
}

func (m *MockDocumentRepository) Toggle() bool {
	m.isOnline = !m.isOnline
	return m.isOnline
}
