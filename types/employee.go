package types

type Employee struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

func NewEmployee(name, position string, salary float64) *Employee {
	return &Employee{
		ID:       123,
		Name:     name,
		Position: position,
		Salary:   salary,
	}
}

type CreateEmployeeRequest struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

type CreateEmployeeResponse struct{}
