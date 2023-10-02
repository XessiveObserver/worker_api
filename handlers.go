package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// getEmployees retrieves all employees in the database and returns in JSON
func getEmployees(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := DB.Query("SELECT * FROM employees")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Salary, &employee.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, employee)
	}
	json.NewEncoder(w).Encode(employees)
}

// getEmployee retrieves a specific employee from the database by ID and returns in JSON
func getEmployee(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(ps.ByName("id"))

	var employee Employee
	err := DB.QueryRow("SELECT * FROM employees WHERE id = $1", id).Scan(&employee.ID, &employee.Name, &employee.Salary, &employee.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Employee not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	json.NewEncoder(w).Encode(employee)
}

// createEmployee insterts a new employee into the database and returns emplyee as JSON
func createEmployee(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "applicatio/json")

	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)

	// Insert new employee into the database
	err := DB.QueryRow("INSERT INTO employees (name, salary, address) VALUES ($1, $2, $3) RETURNING id", employee.Name, employee.Salary, employee.Address).Scan(&employee.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employee)
}

// updateEmployee updates an employee in the database by ID and returns the updated employee in JSON
func updateEmployee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(ps.ByName("id"))

	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)

	// Update user in the database
	_, err := DB.Exec("UPDATE employees SET name = $1, salary = $2, address = $3 WHERE id = $4", employee.Name, employee.Salary, employee.Address, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employee.ID = id
	json.NewEncoder(w).Encode(employee)
}

// deleteEmployee deletes a employee from the database by ID and returns the deleted employee as JSON
func deleteEmployee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(ps.ByName("id"))

	var employee Employee
	err := DB.QueryRow("DELETE FROM employees WHERE id = $1 RETURNING id, name, salary, address", id).Scan(&employee.ID, &employee.Name, &employee.Salary, &employee.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Employee not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employee)
}
