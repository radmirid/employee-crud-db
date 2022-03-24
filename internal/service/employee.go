package service

import (
	"context"
	"time"

	"github.com/radmirid/employee-crud-db/internal/domain"
)

type EmployeesRepository interface {
	Create(ctx context.Context, employee domain.Employee) error
	GetByID(ctx context.Context, id int64) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateEmployeeInput) error
}

type Employees struct {
	repo EmployeesRepository
}

func NewEmployees(repo EmployeesRepository) *Employees {
	return &Employees{
		repo: repo,
	}
}

func (b *Employees) Create(ctx context.Context, employee domain.Employee) error {
	if employee.Birthday.IsZero() {
		employee.Birthday = time.Now()
	}

	return b.repo.Create(ctx, employee)
}

func (b *Employees) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *Employees) GetAll(ctx context.Context) ([]domain.Employee, error) {
	return b.repo.GetAll(ctx)
}

func (b *Employees) Delete(ctx context.Context, id int64) error {
	return b.repo.Delete(ctx, id)
}

func (b *Employees) Update(ctx context.Context, id int64, inp domain.UpdateEmployeeInput) error {
	return b.repo.Update(ctx, id, inp)
}
