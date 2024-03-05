package main

import (
	"fmt"
	"goyav/internal/adapter/web"
	"goyav/internal/core/port"
	"goyav/internal/service"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	var (
		byteRepo      port.ByteRepository
		docRepo       port.DocumentRepository
		analyzer      port.AntivirusAnalyzer
		host          string
		port          int64
		maxUploadSize uint64
		uploadTimeout uint64
		version       string
		information   string
		err           error
	)

	// Setup application configurations
	if err = setup(&host, &port, &maxUploadSize, &uploadTimeout, &version, &information, &byteRepo, &docRepo, &analyzer); err != nil {
		slog.Error("GoyAV failed to setup", "error", err.Error())
		os.Exit(1)
	}

	service, err := service.New(byteRepo, docRepo, analyzer, version, information)
	if err != nil {
		slog.Error("GoyAV failed to initiate the serive", "error", err.Error())
		os.Exit(1)
	}

	// Setting up HTTP server
	mux := web.NewDocumentMux(service, maxUploadSize)
	server := http.Server{
		ReadTimeout: time.Duration(uploadTimeout) * time.Second,
		Addr:        fmt.Sprintf("%v:%v", host, port),
		Handler:     mux,
	}

	// Starting HTTP server
	if err = server.ListenAndServe(); err != nil {
		slog.Error("GoyAV failed to start", "error", err.Error())
		os.Exit(1)
	}
}
