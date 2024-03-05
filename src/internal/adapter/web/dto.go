package web

import (
	"goyav/internal/core/domain"
	"html"
	"time"
)

type DocumentDTO struct {
	ID         string `json:"id"`
	Hash       string `json:"hash"`
	HashAlgo   string `json:"hash_algo"`
	Tag        string `json:"tag"`
	Status     string `json:"analyse_status"`
	AnalyzedAt string `json:"analyzed_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

func NewDocumentDTO(d *domain.Document) *DocumentDTO {
	var (
		status     string
		analyzedAt string
		createdAt  string
		tag        string
	)

	switch d.Status {
	case domain.StatusClean:
		status = "clean"
	case domain.StatusInfected:
		status = "infected"
	default:
		status = "pending"
	}

	if d.Status != domain.StatusPending {
		analyzedAt = d.AnalyzedAt.Format(time.RFC3339)
	}

	createdAt = d.CreatedAt.Format(time.RFC3339)
	tag = html.EscapeString(d.Tag)

	return &DocumentDTO{
		ID:         d.ID,
		Hash:       d.Hash,
		HashAlgo:   "SHA-256",
		Tag:        tag,
		Status:     status,
		CreatedAt:  createdAt,
		AnalyzedAt: analyzedAt,
	}
}
