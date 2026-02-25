package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/dto"
)

// Handler - main HTTP handler struct
type Handler struct {
	services domain.Service
	log      *slog.Logger
}

// NewHandler - constructor for Handler
func NewHandler(services domain.Service, log *slog.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

// CreateDepartment - HTTP handler for creating a new department
func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateDepartment"

	log := h.log.With(slog.String("op", op))
	log.Debug("starting creating department")

	var req dto.CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, h.log, op, err)
		return
	}

	resp, err := h.services.Department().Create(r.Context(), &req)
	if err != nil {
		handleError(w, h.log, op, err)
		return
	}

	log.Info("created department", "name", req.Name)
	renderJSON(w, http.StatusCreated, resp)
}

// GetDepartment - HTTP handler for retrieving a department by ID
func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetDepartment"

	log := h.log.With(slog.String("op", op))
	log.Debug("starting getting department")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, h.log, op, domain.ErrNotFound)
		return
	}

	query := r.URL.Query()
	depth, _ := strconv.Atoi(query.Get("depth"))
	if depth <= 0 {
		depth = 1
	}

	includeEmployees := true
	if query.Get("include_employees") == "false" {
		includeEmployees = false
	}

	req := &dto.GetByIDRequest{
		Depth:            depth,
		IncludeEmployees: includeEmployees,
	}

	resp, err := h.services.Department().GetByID(r.Context(), id, req)
	if err != nil {
		handleError(w, h.log, op, err)
		return
	}

	log.Info("got department", "id", id)
	renderJSON(w, http.StatusOK, resp)
}

// UpdateDepartment - HTTP handler for updating a department by ID
func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateDepartment"

	log := h.log.With(slog.String("op", op))
	log.Debug("starting update department")

	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, h.log, op, err)
		return
	}

	resp, err := h.services.Department().Update(r.Context(), id, &req)
	if err != nil {
		handleError(w, h.log, op, err)
		return
	}

	log.Info("updated department", "id", id)
	renderJSON(w, http.StatusOK, resp)
}

// DeleteDepartment - HTTP handler for deleting a department by ID
func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteDepartment"

	log := h.log.With(slog.String("op", op))
	log.Debug("starting deleting department")

	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	query := r.URL.Query()
	mode := query.Get("mode")
	if mode == "" {
		mode = domain.ModeCascade
	}

	var reassignID *int
	if mode == domain.ModeReassign {
		val, err := strconv.Atoi(query.Get("reassign_to_department_id"))
		if err != nil {
			handleError(w, h.log, op, domain.ErrInvalidReassignToID)
			return
		}
		reassignID = &val
	}

	req := &dto.DeleteDepartmentRequest{
		Mode:         mode,
		ReassignToID: reassignID,
	}

	if err := h.services.Department().Delete(r.Context(), id, req); err != nil {
		handleError(w, h.log, op, err)
		return
	}

	log.Info("deleted department", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

// CreateEmployee - HTTP handler for creating a new employee in a department
func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateEmployee"

	log := h.log.With(slog.String("op", op))
	log.Debug("starting creating employee")

	idStr := r.PathValue("id")
	deptID, _ := strconv.Atoi(idStr)

	var req dto.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, h.log, op, err)
		return
	}

	resp, err := h.services.Employee().Create(r.Context(), deptID, &req)
	if err != nil {
		handleError(w, h.log, op, err)
		return
	}

	log.Info("created employee", "dept_id", deptID)
	renderJSON(w, http.StatusCreated, resp)
}
