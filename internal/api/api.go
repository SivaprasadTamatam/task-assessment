package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tasks/task-assessment/internal/db"
	"github.com/tasks/task-assessment/internal/employee"
)

const (
	DefaultPageSize = 10
)

func ListEmployees(w http.ResponseWriter, r *http.Request) {

	store := r.Context().Value("store").(*db.EmployeeStore)

	pageStr := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	pg, err := strconv.Atoi(pageStr)
	if err != nil || pg <= 0 {
		pg = 1
	}

	pgSize, err := strconv.Atoi(pageSize)
	if err != nil || pgSize <= 0 {
		pgSize = DefaultPageSize
	}
	s := (pg - 1) * pgSize
	e := s + pgSize
	store.Mu.Lock()
	defer store.Mu.Unlock()

	if s >= len(store.Employees) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]employee.Employee{})
		return
	}

	if e > len(store.Employees) {
		e = len(store.Employees)
	}
	employees := store.Employees[s:e]

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
	}
}

func CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var emp employee.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store := r.Context().Value("store").(*db.EmployeeStore)
	emp, err = store.CreateEmployee(emp)

	if err != nil {
		http.Error(w, "Failed to creae employee object", http.StatusInternalServerError)
		return
	}
	// Respond with the created employee in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func GetEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the employee ID from the URL path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	store := r.Context().Value("store").(*db.EmployeeStore)
	emp, err := store.GetEmployeeByID(id)

	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var updatedEmp employee.Employee

	err := json.NewDecoder(r.Body).Decode(&updatedEmp)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store := r.Context().Value("store").(*db.EmployeeStore)
	updatedEmp, err = store.UpdateEmployee(updatedEmp.ID, updatedEmp)

	if err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmp)
}

func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	store := r.Context().Value("store").(*db.EmployeeStore)
	err = store.DeleteEmployee(id)

	if err != nil {
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Employee with ID %d deleted successfully", id)
}

func NewAPI(store *db.EmployeeStore) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/employees", ListEmployees).Methods("GET")
	r.HandleFunc("/employee", CreateEmployeeHandler).Methods("POST")
	r.HandleFunc("/employee/{id}", GetEmployeeByIDHandler).Methods("GET")
	r.HandleFunc("/employee", UpdateEmployeeHandler).Methods("PUT")
	r.HandleFunc("/employee/{id}", DeleteEmployeeHandler).Methods("DELETE")

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "store", store)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	return r
}
