package docrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goyav/internal/core/domain"
	"goyav/internal/core/port"
	"time"

	_ "github.com/lib/pq"
)

type PotgresDocumentRepository struct {
	db *sql.DB
}

var ErrPostgresDocumentRepository = errors.New("PostgresDocumentRepository")

func NewPotgresDocumentRepository(host string, port uint64, dbname string, schema string, user string, password string, useSSL bool) (*PotgresDocumentRepository, error) {
	var ssl = "disable"
	if useSSL {
		ssl = "enable"
	}

	connInfo := fmt.Sprintf("host=%v port=%v dbname=%v search_path=%v sslmode=%v user=%v password=%v", host, port, dbname, schema, ssl, user, password)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to create document repository:  %v", ErrPostgresDocumentRepository, err)
	}
	return &PotgresDocumentRepository{
		db: db,
	}, nil
}

// Save adds a new document to the repository and returns an error if the document already exists or
// if there is an issue during the save operation.
func (r PotgresDocumentRepository) Save(ctx context.Context, doc *domain.Document) error {
	q := "INSERT INTO documents (document_id, hash, tag, status, analyzed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.ExecContext(ctx, q, doc.ID, doc.Hash, doc.Tag, doc.Status, doc.AnalyzedAt, doc.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: %w: %v: document=%#v", ErrPostgresDocumentRepository, port.ErrSaveDocumentFailed, err, doc)
	}
	return nil
}

// Get retrieves a document by its ID and returns an error if not found or if there is an issue with the ID.
func (r PotgresDocumentRepository) Get(ctx context.Context, ID string) (*domain.Document, error) {
	q := "SELECT document_id, hash, tag, status, analyzed_at, created_at FROM documents WHERE document_id = $1"
	doc := new(domain.Document)
	err := r.db.QueryRowContext(ctx, q, ID).Scan(
		&doc.ID,
		&doc.Hash,
		&doc.Tag,
		&doc.Status,
		&doc.AnalyzedAt,
		&doc.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrDocumentNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrGetDocumentFailed, err)
	}
	return doc, nil
}

// GetByHash retrieves a document by its hash and returns an error if not found or if there is an issue with the hash.
func (r PotgresDocumentRepository) GetByHash(ctx context.Context, hash string) (*domain.Document, error) {
	q := "SELECT document_id, hash, tag, status, analyzed_at, created_at FROM documents WHERE hash = $1"
	doc := new(domain.Document)
	err := r.db.QueryRowContext(ctx, q, hash).Scan(
		&doc.ID,
		&doc.Hash,
		&doc.Tag,
		&doc.Status,
		&doc.AnalyzedAt,
		&doc.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w.GetByHash: %w", ErrPostgresDocumentRepository, port.ErrDocumentNotFound)
		}

		return nil, fmt.Errorf("%w.GetByHash: %w: %v", ErrPostgresDocumentRepository, port.ErrGetDocumentFailed, err)
	}
	return doc, nil
}

// Delete removes a document from the repository by its ID and returns an error if not found or during deletion.
func (r PotgresDocumentRepository) Delete(ctx context.Context, ID string) error {
	q := "DELETE FROM documents WHERE document_id = $1"
	_, err := r.db.ExecContext(ctx, q, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrDeleteDocumentFailed, err)
		}

		return fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrDeleteDocumentFailed, err)
	}
	return nil
}

// UpdateStatus updates a document's analysis status and date, returning an error for nonexistent documents,
// invalid status, or update issues.
func (r PotgresDocumentRepository) UpdateStatus(ctx context.Context, ID string, status domain.AnalysisStatus, analyzedAt time.Time) error {
	q := "UPDATE documents SET status = $2, analyzed_at = $3 WHERE document_id = $1"
	_, err := r.db.ExecContext(ctx, q, ID, status, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrUpdateStatusFailed, err)
		}

		return fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrUpdateStatusFailed, err)
	}
	return nil
}

// Ping checks the repository's availability or health status.
func (r PotgresDocumentRepository) Ping() error {
	if err := r.db.Ping(); err != nil {
		return fmt.Errorf("%w: %w: %v", ErrPostgresDocumentRepository, port.ErrDocumentRepositoryUnavailable, err)
	}
	return nil
}
