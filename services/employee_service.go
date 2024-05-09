package services

import (
	"database/sql"

	"github.com/kevalsabhani/go-assignment/types"
)

type EmployeeService struct {
	DB *sql.DB
}

func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		DB: db,
	}
}
func (service *EmployeeService) Create(reqParams *types.CreateEmployeeRequest) (*types.Employee, error) {
	var employee types.Employee
	err := service.DB.QueryRow(
		"INSERT INTO employees(name, position, salary) VALUES($1, $2, $3) RETURNING id",
		reqParams.Name, reqParams.Position, reqParams.Salary).Scan(&employee.ID)

	if err != nil {
		return nil, err
	}
	employee.Name = reqParams.Name
	employee.Position = reqParams.Position
	employee.Salary = reqParams.Salary
	return &employee, nil
}

func (service *EmployeeService) GetById(id int) (*types.Employee, error) {
	var employee types.Employee
	row := service.DB.QueryRow("SELECT id, name, position, salary FROM employees WHERE id=$1", id)
	if err := row.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary); err != nil {
		return nil, err
	}
	return &employee, nil
}

func (service *EmployeeService) Get(page int, size int) ([]*types.Employee, error) {
	rows, err := service.DB.Query(
		"SELECT id, name, position, salary FROM employees LIMIT $1 OFFSET $2",
		size, (page-1)*size)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	employees := []*types.Employee{}

	for rows.Next() {
		var emp types.Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary); err != nil {
			return nil, err
		}
		employees = append(employees, &emp)
	}

	return employees, nil
}

func (service *EmployeeService) Update(id int, reqParams *types.CreateEmployeeRequest) error {
	_, err :=
		service.DB.Exec("UPDATE employees SET name=$1, position=$2, salary=$3 WHERE id=$4",
			reqParams.Name, reqParams.Position, reqParams.Salary, id)

	return err
}

func (service *EmployeeService) Delete(id int) error {
	_, err := service.DB.Exec("DELETE FROM employees WHERE id=$1", id)
	return err
}
