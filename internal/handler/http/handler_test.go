package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/dto"
)

// MOCKS

type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) Create(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) GetByID(ctx context.Context, id int, req *dto.GetByIDRequest) (*dto.DepartmentResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) Update(ctx context.Context, id int, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) Delete(ctx context.Context, id int, req *dto.DeleteDepartmentRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) Create(ctx context.Context, deptID int, req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	args := m.Called(ctx, deptID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EmployeeResponse), args.Error(1)
}

type MockService struct {
	mock.Mock
	dept domain.DepartmentService
	emp  domain.EmployeeService
}

func (m *MockService) Department() domain.DepartmentService { return m.dept }
func (m *MockService) Employee() domain.EmployeeService     { return m.emp }

func setupTest(t *testing.T) (*MockDepartmentService, *MockEmployeeService, *http.ServeMux) {
	mockDept := new(MockDepartmentService)
	mockEmp := new(MockEmployeeService)

	// mock
	mockSrv := &MockService{
		dept: mockDept,
		emp:  mockEmp,
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewHandler(mockSrv, logger)
	mux := NewRouter(h)

	return mockDept, mockEmp, mux
}

// TESTS

func TestHandler_CreateDepartment(t *testing.T) {
	mockDept, _, mux := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateDepartmentRequest{Name: "IT"}
		resp := &dto.DepartmentResponse{ID: 1, Name: "IT"}

		mockDept.On("Create", mock.Anything, mock.MatchedBy(func(r *dto.CreateDepartmentRequest) bool {
			return r.Name == "IT"
		})).Return(resp, nil).Once()

		body, _ := json.Marshal(req)
		r := httptest.NewRequest("POST", "/departments", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestHandler_GetDepartment(t *testing.T) {
	mockDept, _, mux := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		deptID := 1
		// Хендлер по умолчанию ставит depth=1 и include_employees=true
		expectedReq := &dto.GetByIDRequest{Depth: 1, IncludeEmployees: true}
		resp := &dto.DepartmentResponse{ID: 1, Name: "IT"}

		mockDept.On("GetByID", mock.Anything, deptID, expectedReq).Return(resp, nil).Once()

		r := httptest.NewRequest("GET", "/departments/1", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockDept.On("GetByID", mock.Anything, 99, mock.Anything).Return(nil, domain.ErrDepartmentNotFound).Once()

		r := httptest.NewRequest("GET", "/departments/99", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_DeleteDepartment(t *testing.T) {
	mockDept, _, mux := setupTest(t)

	t.Run("Success Cascade", func(t *testing.T) {
		// По умолчанию mode=cascade
		req := &dto.DeleteDepartmentRequest{Mode: "cascade", ReassignToID: nil}

		mockDept.On("Delete", mock.Anything, 1, req).Return(nil).Once()

		r := httptest.NewRequest("DELETE", "/departments/1", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
