package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"goyav/internal/core/port"
	"log"
	"log/slog"
	"net/http"
)

func (d *DocumentMux) root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ping", http.StatusPermanentRedirect)
}

func (d *DocumentMux) getDocumentByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, r)
		return
	}
	om := &ObjectMessage{}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "please provide a document ID", om)
		return
	}
	doc, err := d.service.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, port.ErrDocumentNotFound) {
			writeError(w, http.StatusNotFound, "document not found", om)
		} else {
			writeError(w, http.StatusInternalServerError, "an error occured", om)
			log.Printf("getDocumentHandler: %v", err.Error())
		}
		return
	} else {
		om.Message = "document found"
		om.Document = NewDocumentDTO(doc)
	}
	b, err := json.Marshal(om)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "an error occured", om)
		log.Printf("getDocumentHandler: error while marshalling : %s", err.Error())
		return
	}
	w.Write(b)
}

func (d *DocumentMux) postDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, r)
		return
	}

	var (
		om               = &ObjectMessage{}
		reqSizeLim int64 = int64(d.maxUploadSize) + (2 << 10)
	)

	r.Body = http.MaxBytesReader(w, r.Body, reqSizeLim)
	defer r.Body.Close()
	if err := r.ParseMultipartForm(reqSizeLim); err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("uploaded data exceeds the maximum allowed size : %v Bytes.", d.maxUploadSize), om)
		slog.Debug(fmt.Sprintf("handler.postDocumentHandler: %v", om.Message), "error", err.Error())
		return
	}

	tag := r.FormValue("tag")

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to upload file", om)
		slog.Error("handler.postDocumentHandler: "+om.Message, "msg", err.Error())
		return
	}
	defer file.Close()

	if header.Size == 0 {
		writeError(w, http.StatusBadRequest, "the file to upload is empty", om)
		return
	}

	if tag == "" {
		tag = header.Filename
	}
	ID, err := d.service.Upload(r.Context(), file, header.Size, tag)
	if err != nil {
		if errors.Is(err, port.ErrDocumentAlreadyExists) {
			om.Message = "document already exists."
		} else {
			writeError(w, http.StatusInternalServerError, "an error occured while uploading", om)
			slog.Error("handler.postDocumentHandler: "+om.Message, "msg", err.Error())
		}
	} else {
		om.Message = "document uploaded successfully."
		w.WriteHeader(http.StatusCreated)
	}
	if ID != "" {
		om.ID = ID
	}
	writeJson(w, om)
}

func (d *DocumentMux) ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, r)
		return
	}
	om := &ObjectMessage{
		Information: d.service.Information(),
		Version:     d.service.Version(),
	}
	if err := d.service.Ping(); err != nil {
		writeError(w, http.StatusServiceUnavailable, "service unavailable", om)
	} else {
		w.WriteHeader(http.StatusOK)
		om.Message = "PONG : everything is good"
	}
	writeJson(w, om)
}
