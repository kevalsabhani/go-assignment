package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kevalsabhani/go-assignment/services"
	"github.com/kevalsabhani/go-assignment/types"
	"github.com/kevalsabhani/go-assignment/utils"
)

type EmployeeHandler struct {
	Service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		Service: service,
	}
}

// Create employee
func (handler *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var empParams types.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&empParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	emp, err := handler.Service.Create(&empParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, emp)
}

// Get employee By ID
func (handler *EmployeeHandler) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	employee, err := handler.Service.GetById(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "employee not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, employee)

}

// Get employees
func (handler *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	var page, size int
	var err error
	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid page number")
			return
		}
	}
	if r.URL.Query().Get("size") == "" {
		size = 10
	} else {
		size, err = strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid size")
			return
		}
	}
	if page < 1 || size < 1 {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid page/size")
	}

	employees, err := handler.Service.Get(page, size)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, employees)

}

// Update employee
func (handler *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	var empParams types.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&empParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := handler.Service.Update(id, &empParams); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Delete employee
func (handler *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	if err := handler.Service.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
