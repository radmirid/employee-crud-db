package domain

import (
	"time"
)

type Employee struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Birthday time.Time `json:"birthday"`
	Utility  int       `json:"utility"`
}

type UpdateEmployeeInput struct {
	Name     *string    `json:"name"`
	Surname  *string    `json:"surname"`
	Birthday *time.Time `json:"birthday"`
	Utility  *int       `json:"utility"`
}
