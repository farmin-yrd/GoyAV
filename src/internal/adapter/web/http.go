package web

import (
	"goyav/internal/core/port"
	"net/http"
)

// Default upload size limit in bytes : 1 Mib
const DefaultMaxUploadSize int64 = 1 << 20

type DocumentMux struct {
	*http.ServeMux
	service       port.DocumentService
	maxUploadSize uint64
}

func NewDocumentMux(s port.DocumentService, n uint64) *DocumentMux {
	d := &DocumentMux{
		ServeMux:      http.NewServeMux(),
		maxUploadSize: n << 20,
		service:       s,
	}
	d.setup()
	return d
}
