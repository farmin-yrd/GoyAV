package service

import (
	"context"
	"errors"
	"fmt"
	"goyav/internal/core/domain"
	"goyav/internal/core/port"
	"goyav/pkg/helper"
	"io"
	"log/slog"
	"time"
)

// Service manages file uploads and antivirus analysis operations.
type Service struct {
	ByteRepository     port.ByteRepository
	DocumentRepository port.DocumentRepository
	AvAnalyzer         port.AntivirusAnalyzer

	// AnalysisAttempts specifies the number of times an analysis is attempted in case of failure.
	AnalysisAttempts uint8

	// semaphore is used to control concurrent access to resources.
	semaphore chan struct{}

	// version is the current version of the service
	version string

	// information about the service
	information string
}

const (
	DefaultSemaphoreCapacity = 128
	DefaultAnalysisAttempts  = 3
)

var ErrNilDependency = errors.New("Service : nil dependency")

// New creates a new Service instance with the provided dependencies and information about the service.
func New(b port.ByteRepository, d port.DocumentRepository, a port.AntivirusAnalyzer, ver, info string) (*Service, error) {
	if b == nil || d == nil || a == nil {
		return nil, fmt.Errorf("%w: missing repositories or analyzer", ErrNilDependency)
	}

	if err := ping(b, d, a); err != nil {
		return nil, fmt.Errorf("service: unable to create: %w", err)
	}

	return &Service{
		ByteRepository:     b,
		DocumentRepository: d,
		AvAnalyzer:         a,
		AnalysisAttempts:   DefaultAnalysisAttempts,
		semaphore:          make(chan struct{}, DefaultSemaphoreCapacity),
		version:            ver,
		information:        info,
	}, nil

}

// SetCapacity sets the semaphore capacity.
func (s *Service) SetCapacity(n uint) {
	s.semaphore = make(chan struct{}, n)
}

// Version returns the current version of the service.
func (s *Service) Version() string {
	return s.version
}

// Information returns the information about the service
func (s *Service) Information() string {
	return s.information
}

// Upload uploads a file, triggers asynchronous analysis, and returns its unique identifier.
// If a document with the same hash exists, its previous analysis result is reused.
// The "tag" parameter represents a user-defined tag associated with the document.
func (s *Service) Upload(ctx context.Context, data io.ReadSeeker, size int64, tag string) (ID string, err error) {
	// Verify if a document with the same hash exists
	hash, err := helper.MakeHash(io.LimitReader(data, size))
	if err != nil {
		return "", fmt.Errorf("service: failed to calculate the hash: %w", err)
	}

	if doc, err := s.DocumentRepository.GetByHash(ctx, hash); err == nil {
		return doc.ID, fmt.Errorf("service: %w", port.ErrDocumentAlreadyExists)
	} else if !errors.Is(err, port.ErrDocumentNotFound) {
		return "", fmt.Errorf("service: %w: failed to get document by its hash=%v: error: %v", port.ErrServiceUploadFailed, hash, err)
	}

	// Save the bytes of the document
	if err = seek(data); err != nil {
		return "", fmt.Errorf("service: %w", err)
	}
	ID, err = helper.MakeID(io.LimitReader(data, size))
	if err != nil {
		return "", fmt.Errorf("service: failed to create a document ID: %w", err)
	}
	if err != seek(data) {
		return "", fmt.Errorf("service: %w", err)
	}
	if err = s.ByteRepository.Save(ctx, io.LimitReader(data, size), size, ID); err != nil {
		return "", fmt.Errorf("service: %w: %w: id=%v", port.ErrServiceUploadFailed, err, ID)
	}

	// Create a new document instance with the file's ID, hash, and truncated tag.
	tag = helper.TruncateTag(tag)
	newDoc := domain.NewDocument(ID, hash, tag)

	// Attempt to retrieve an existing document with the same hash.
	existing, _ := s.DocumentRepository.GetByHash(ctx, hash)

	// If an existing document is found and its analyse status is not pending, update the new document with the existing one.
	if existing != nil && existing.Status != domain.StatusPending {
		newDoc.Status = existing.Status
		newDoc.AnalyzedAt = existing.AnalyzedAt
	} else {
		// Otherwise, trigger the file analysis asynchronously
		go s.asyncAnalyze(ID)
	}

	// Create the new document in the document repository.
	err = s.DocumentRepository.Save(ctx, newDoc)
	if err != nil {
		return "", fmt.Errorf("Service: %w: %w", port.ErrServiceUploadFailed, err)
	}
	return ID, nil
}

// GetDocument retrieves the current status of a document by its ID.
func (s *Service) GetDocument(ctx context.Context, ID string) (*domain.Document, error) {
	document, err := s.DocumentRepository.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: id=%s", port.ErrServiceGetDocumentFailed, err, ID)
	}
	return document, nil
}

func (s *Service) Ping() error {
	return ping(s.ByteRepository, s.DocumentRepository, s.AvAnalyzer)
}

const asyncAnalyseErrorMsg = "service : async analysis error"

// asyncAnalyze performs the analysis of the data asynchronously with retry attempts
func (s *Service) asyncAnalyze(ID string) {
	s.semaphore <- struct{}{}
	go func() {
		defer func() {
			<-s.semaphore
		}()

		ctx := context.Background()

		// Retrieve and defer close the data stream
		r, err := s.ByteRepository.Get(ctx, ID)
		if err != nil {
			slog.Error(asyncAnalyseErrorMsg, "error", err, "ID", ID)
			return
		}
		defer r.Close()

		// Attempt to analyze with retries
		if err := s.attemptAnalysis(ctx, r, ID); err != nil {
			slog.Error(asyncAnalyseErrorMsg, "error", err, "ID", ID)
		}
		slog.Debug("analyse completed", "ID", ID)

	}()
}

// attemptAnalysis tries to analyze the data with retries.
func (s *Service) attemptAnalysis(ctx context.Context, r io.Reader, ID string) error {
	var status domain.AnalysisStatus
	for i := uint8(1); i <= s.AnalysisAttempts; i++ {
		var err error
		if status, err = s.AvAnalyzer.Analyze(ctx, r); err == nil {
			if err = s.DocumentRepository.UpdateStatus(ctx, ID, status, time.Now()); err == nil {
				return s.ByteRepository.Delete(ctx, ID)
			}
			return err
		}
		time.Sleep(time.Second * 30 * time.Duration(i))
	}
	return fmt.Errorf("analysis failed after %d attempts", s.AnalysisAttempts)
}

func seek(r io.ReadSeeker) error {
	_, err := r.Seek(0, io.SeekStart)
	return err
}

func ping(b port.ByteRepository, d port.DocumentRepository, a port.AntivirusAnalyzer) error {
	return errors.Join(b.Ping(), d.Ping(), a.Ping())
}
