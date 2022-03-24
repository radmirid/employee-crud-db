package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/radmirid/employee-crud-db/internal/domain"
)

func (h *Handler) getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		logError("getEmployeeByID", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	employee, err := h.employeesService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorEmployeeNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logError("getEmployeeByID", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(employee)
	if err != nil {
		logError("getEmployeeByID", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("createEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var employee domain.Employee
	if err = json.Unmarshal(reqBytes, &employee); err != nil {
		logError("createEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employeesService.Create(r.Context(), employee)
	if err != nil {
		logError("createEmployee", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		logError("deleteEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employeesService.Delete(r.Context(), id)
	if err != nil {
		logError("deleteEmployee", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.employeesService.GetAll(r.Context())
	if err != nil {
		logError("getAllEmployees", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(employees)
	if err != nil {
		logError("getAllEmployees", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		logError("updateEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("updateEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateEmployeeInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logError("updateEmployee", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employeesService.Update(r.Context(), id, inp)
	if err != nil {
		logError("updateEmployee", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
