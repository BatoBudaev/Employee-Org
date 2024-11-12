package models

type Employee struct {
	ID             int
	FirstName      string
	LastName       string
	OrganizationID int
	Photo          []byte
}
