package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/radmirid/employee-crud-db/internal/domain"

	"github.com/gorilla/mux"
)

type Employees interface {
	Create(ctx context.Context, employee domain.Employee) error
	GetByID(ctx context.Context, id int64) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateEmployeeInput) error
}

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	employeesService Employees
	usersService     User
}

func NewHandler(employees Employees, users User) *Handler {
	return &Handler{
		employeesService: employees,
		usersService:     users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(logger)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)

	}

	employees := r.PathPrefix("/employees").Subrouter()
	{
		employees.Use(h.authorizer)

		employees.HandleFunc("", h.createEmployee).Methods(http.MethodPost)
		employees.HandleFunc("", h.getAllEmployees).Methods(http.MethodGet)
		employees.HandleFunc("/{id:[0-9]+}", h.getEmployeeByID).Methods(http.MethodGet)
		employees.HandleFunc("/{id:[0-9]+}", h.deleteEmployee).Methods(http.MethodDelete)
		employees.HandleFunc("/{id:[0-9]+}", h.updateEmployee).Methods(http.MethodPut)
	}

	return r
}

func getIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("error: zero id")
	}

	return id, nil
}
