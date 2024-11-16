package models

type Organization struct {
	ID      int
	Name    string
	Address string
}

type Employee struct {
	ID               int
	FirstName        string
	LastName         string
	OrganizationID   int
	Photo            []byte
	OrganizationName string
}
