package db

import (
	"errors"
	"sync"

	"github.com/tasks/task-assessment/internal/employee"
)

type EmployeeStore struct {
	Employees []employee.Employee
	Mu        sync.RWMutex
	NextID    int
}

func NewEmployeeStore() *EmployeeStore {
	return &EmployeeStore{
		Employees: []employee.Employee{},
		NextID:    1,
	}
}

func (store *EmployeeStore) CreateEmployee(emp employee.Employee) (employee.Employee, error) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	emp.ID = store.NextID
	store.NextID++
	store.Employees = append(store.Employees, emp)

	return emp, nil
}

func (store *EmployeeStore) GetEmployeeByID(id int) (employee.Employee, error) {
	store.Mu.RLock()
	defer store.Mu.RUnlock()

	for _, emp := range store.Employees {
		if emp.ID == id {
			return emp, nil
		}
	}

	return employee.Employee{}, errors.New("employee not found")
}

func (store *EmployeeStore) UpdateEmployee(id int, updateEmp employee.Employee) (employee.Employee, error) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	for i, emp := range store.Employees {
		if emp.ID == id {
			store.Employees[i] = updateEmp
			return updateEmp, nil
		}
	}

	return employee.Employee{}, errors.New("employee not found")
}

func (store *EmployeeStore) DeleteEmployee(id int) error {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	for i, emp := range store.Employees {
		if emp.ID == id {
			store.Employees = append(store.Employees[:i], store.Employees[i+1:]...)
			return nil
		}
	}
	return errors.New("employee not found")
}
