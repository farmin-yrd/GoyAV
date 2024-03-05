package antivirus

import (
	"bytes"
	"context"
	"fmt"
	"goyav/internal/core/domain"
	"goyav/internal/core/port"
	"io"
	"time"
)

// MockAntivirusAnalyzer is a mock implementation of the AntivirusAnalyzer interface.
// It uses the EICAR test byte slice to simulate virus detection.
type MockAntivirusAnalyzer struct {
	Online  bool          // Indicates whether the mock analyzer is "online" or "offline"
	Timeout time.Duration // Timeout in seconds
}

// NewMock creates a new instance of MockAntivirusAnalyzer.
func NewMock() *MockAntivirusAnalyzer {
	return &MockAntivirusAnalyzer{
		Online:  true,
		Timeout: 60 * time.Second,
	}
}

// Analyze performs a mock antivirus analysis on the byte content of a document.
func (m *MockAntivirusAnalyzer) Analyze(ctx context.Context, r io.Reader) (domain.AnalysisStatus, error) {
	var (
		status  domain.AnalysisStatus = domain.StatusPending
		err     error
		jobDone = make(chan struct{})
	)

	go func() {
		defer close(jobDone)
		var b []byte
		b, err = io.ReadAll(r)
		if err == nil {
			if bytes.Contains(b, port.EICAR) {
				status = domain.StatusInfected
			} else {
				status = domain.StatusClean
			}
		}
		time.Sleep(time.Second * 2)
		jobDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		// If the context is cancelled or times out, return a status indicating the operation didn't complete
		return domain.StatusPending, fmt.Errorf("MockAntivirusAnalyzer: %w: %v", port.ErrAntivirusAnalysisFailed, ctx.Err())
	case <-jobDone:
		// Return the analysis status if the job is done
		return status, err
	}
}

// Ping simulates a connectivity check to the antivirus service.
func (m *MockAntivirusAnalyzer) Ping() error {
	if !m.Online {
		return fmt.Errorf("MockAntivirusAnalyzer: %w", port.ErrAntivirusAnalyserUnavailable)
	}
	return nil
}

func (m *MockAntivirusAnalyzer) TimeoutValue() uint64 {
	return uint64(m.Timeout.Seconds())
}
