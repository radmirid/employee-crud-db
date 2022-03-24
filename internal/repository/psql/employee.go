package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/radmirid/employee-crud-db/internal/domain"
)

type Employees struct {
	db *sql.DB
}

func NewEmployees(db *sql.DB) *Employees {
	return &Employees{db}
}

func (b *Employees) Create(ctx context.Context, employee domain.Employee) error {
	_, err := b.db.Exec("INSERT INTO employees (name, surname, birthday, utility) values ($1, $2, $3, $4)",
		employee.Name, employee.Surname, employee.Birthday, employee.Utility)

	return err
}

func (b *Employees) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	var employee domain.Employee
	err := b.db.QueryRow("SELECT id, name, surname, birthday, utility FROM employees WHERE id=$1", id).
		Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.Birthday, &employee.Utility)
	if err == sql.ErrNoRows {
		return employee, domain.ErrorEmployeeNotFound
	}

	return employee, err
}

func (b *Employees) GetAll(ctx context.Context) ([]domain.Employee, error) {
	rows, err := b.db.Query("SELECT id, name, surname, birthday, utility FROM employees")
	if err != nil {
		return nil, err
	}

	employees := make([]domain.Employee, 0)
	for rows.Next() {
		var employee domain.Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.Birthday, &employee.Utility); err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, rows.Err()
}

func (b *Employees) Delete(ctx context.Context, id int64) error {
	_, err := b.db.Exec("DELETE FROM employees WHERE id=$1", id)

	return err
}

func (b *Employees) Update(ctx context.Context, id int64, inp domain.UpdateEmployeeInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *inp.Name)
		argID++
	}

	if inp.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argID))
		args = append(args, *inp.Surname)
		argID++
	}

	if inp.Birthday != nil {
		setValues = append(setValues, fmt.Sprintf("birthday=$%d", argID))
		args = append(args, *inp.Birthday)
		argID++
	}

	if inp.Utility != nil {
		setValues = append(setValues, fmt.Sprintf("utility=$%d", argID))
		args = append(args, *inp.Utility)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE employees SET %s WHERE id=$%d", setQuery, argID)
	args = append(args, id)

	_, err := b.db.Exec(query, args...)
	return err
}
