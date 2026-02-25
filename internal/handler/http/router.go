package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tmozzze/org_struct_api/docs"
)

// NewRouter - creates and returns a new HTTP router with all routes registered.
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	// Swagger UI
	// http://localhost:8080/swagger/index.html
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Departments
	mux.HandleFunc("POST /departments", h.CreateDepartment)
	mux.HandleFunc("GET /departments/{id}", h.GetDepartment)
	mux.HandleFunc("PATCH /departments/{id}", h.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", h.DeleteDepartment)

	// Employees
	mux.HandleFunc("POST /departments/{id}/employees", h.CreateEmployee)

	return mux
}
