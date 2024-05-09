package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevalsabhani/go-assignment/db"
	"github.com/kevalsabhani/go-assignment/handlers"
	"github.com/kevalsabhani/go-assignment/services"
	"github.com/kevalsabhani/go-assignment/types"
)

var testDB *sql.DB
var r *mux.Router

func TestMain(m *testing.M) {
	// Load environment variables
	godotenv.Load()

	// Connect DB
	testDB = db.ConnectDB()
	r = mux.NewRouter()

	employeeHandler := handlers.NewEmployeeHandler(services.NewEmployeeService(testDB))
	handlers.SetupRoutes(r, employeeHandler)

	os.Exit(m.Run())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addEmployees(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		testDB.Exec("INSERT INTO employees(name, position, salary) VALUES($1, $2, $3)", "Employee "+strconv.Itoa(i), "SE "+strconv.Itoa(i), (i+1.0)*1000)
	}
}

func TestEmptyTable(t *testing.T) {
	db.ClearEmployeeTable(testDB)

	req, _ := http.NewRequest("GET", "/employees", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentEmployee(t *testing.T) {

	req, _ := http.NewRequest("GET", "/employees/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "employee not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'employee not found'. Got '%s'", m["error"])
	}
}

func TestCreateEmployee(t *testing.T) {
	var jsonStr = []byte(`{"name":"test employee", "position": "software engineer", "salary": 23000.22}`)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test employee" {
		t.Errorf("Expected employee name to be 'test employee'. Got '%v'", m["name"])
	}

	if m["position"] != "software engineer" {
		t.Errorf("Expected employee position to be 'software engineer'. Got '%v'", m["position"])
	}

	if m["salary"] != 23000.22 {
		t.Errorf("Expected employee salary to be '23000.22'. Got '%v'", m["salary"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected employee ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetEmployeeById(t *testing.T) {
	db.ClearEmployeeTable(testDB)
	addEmployees(1)

	req, _ := http.NewRequest("GET", "/employees/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetEmployees(t *testing.T) {
	db.ClearEmployeeTable(testDB)
	expectedCount := 3
	addEmployees(expectedCount)

	req, _ := http.NewRequest("GET", "/employees", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var employees []*types.Employee
	json.Unmarshal(response.Body.Bytes(), &employees)

	count := len(employees)
	if count != expectedCount {
		t.Errorf("Expected %v employees, but got %v", expectedCount, count)
	}
}

func TestUpdateEmployee(t *testing.T) {

	db.ClearEmployeeTable(testDB)
	addEmployees(1)

	req, _ := http.NewRequest("GET", "/employees/1", nil)
	response := executeRequest(req)
	var originalEmployee map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEmployee)

	var jsonStr = []byte(`{"name":"test employee - updated name", "position": "SE 1", "salary": 25000.40}`)
	req, _ = http.NewRequest("PUT", "/employees/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["result"] != "success" {
		t.Errorf("Expected result to be success")
	}
}

func TestDeleteEmployee(t *testing.T) {
	db.ClearEmployeeTable(testDB)
	addEmployees(1)

	req, _ := http.NewRequest("GET", "/employees/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/employees/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/employees/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
