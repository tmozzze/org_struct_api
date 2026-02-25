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

// CreateDepartment godoc
// @Summary Create department
// @Description Create department
// @Tags departments
// @Accept json
// @Produce json
// @Param input body dto.CreateDepartmentRequest true "Department data"
// @Success 201 {object} dto.DepartmentResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Router /departments [post]
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

// GetDepartment godoc
// @Summary Get department with details
// @Description Return department with children and employees
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param depth query int false "Tree depth (1-5)" default(1)
// @Param include_employees query bool false "With employees" default(true)
// @Success 200 {object} dto.DepartmentResponse
// @Failure 404 {object} errorResponse
// @Router /departments/{id} [get]
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

// UpdateDepartment godoc
// @Summary Update department
// @Description Update name and parent ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param input body dto.UpdateDepartmentRequest true "New data"
// @Success 200 {object} dto.DepartmentResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Router /departments/{id} [patch]
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

// DeleteDepartment godoc
// @Summary Delete department
// @Description Delete department in cascade mode or reassign mode
// @Tags departments
// @Param id path int true "Department ID"
// @Param mode query string false "Delete mode (cascade|reassign)" Enums(cascade, reassign) default(cascade)
// @Param reassign_to_department_id query int false "New department ID (need for reassign mode)"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /departments/{id} [delete]
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

// CreateEmployee godoc
// @Summary Create employee
// @Description Create employee in department
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param input body dto.CreateEmployeeRequest true "Employee data"
// @Success 201 {object} dto.EmployeeResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /departments/{id}/employees [post]
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
