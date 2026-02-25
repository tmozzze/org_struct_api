package http

import "net/http"

// NewRouter - creates and returns a new HTTP router with all routes registered.
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Departments
	mux.HandleFunc("POST /departments", h.CreateDepartment)
	mux.HandleFunc("GET /departments/{id}", h.GetDepartment)
	mux.HandleFunc("PATCH /departments/{id}", h.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", h.DeleteDepartment)

	// Employees
	mux.HandleFunc("POST /departments/{id}/employees", h.CreateEmployee)

	return mux
}
