package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/tmozzze/org_struct_api/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func handleError(w http.ResponseWriter, log *slog.Logger, op string, err error) {
	log.Error(op, slog.String("err", err.Error()))

	status := http.StatusInternalServerError
	message := "internal server error"

	// Mapping domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrNotFound),
		errors.Is(err, domain.ErrDepartmentNotFound),
		errors.Is(err, domain.ErrParentNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrDuplicateName),
		errors.Is(err, domain.ErrCycleConstraint):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidReassignToID),
		errors.Is(err, domain.ErrLengthConstraint),
		errors.Is(err, domain.ErrEmptyConstraint):
		status = http.StatusBadRequest
		message = err.Error()
	}

	renderJSON(w, status, errorResponse{Error: message})
}
