package handlers

import "github.com/gorilla/mux"

func SetupRoutes(r *mux.Router, h *EmployeeHandler) {
	r.HandleFunc("/employees", h.CreateEmployee).Methods("POST")
	r.HandleFunc("/employees", h.GetEmployees).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", h.GetEmployeeById).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", h.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id:[0-9]+}", h.DeleteEmployee).Methods("DELETE")
}
