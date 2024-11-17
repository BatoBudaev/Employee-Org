package db

import (
	"database/sql"
	"employeeOrgDB/internal/models"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	*sql.DB
}

func InitDB(user, password, dbname string) (*DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8", user, password, dbname)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Подключено к базе данных")

	return &DB{db}, nil
}

func (db *DB) GetEmployess() ([]models.Employee, error) {
	query := `
		SELECT e.id, e.first_name, e.last_name, e.organization_id, o.name
		FROM employees e
		LEFT JOIN organizations o ON e.organization_id = o.id
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.OrganizationID, &emp.OrganizationName); err != nil {
			log.Println("Ошибка чтения строки: ", err)
			continue
		}

		employees = append(employees, emp)
	}

	return employees, nil
}

func (db *DB) GetOrganizations() ([]models.Organization, error) {
	rows, err := db.Query("SELECT id, name, address FROM organizations")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var organizations []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Address); err != nil {
			log.Println("Ошибка чтения строки:", err)
			continue
		}

		organizations = append(organizations, org)
	}

	return organizations, nil
}

func (db *DB) AddEmployee(firstName, lastName string, organizationID int) error {
	_, err := db.Exec("INSERT INTO employees (first_name, last_name, organization_id) VALUES ($1, $2, $3)", firstName, lastName, organizationID)
	return err
}

func (db *DB) AddOrganization(name, address string) error {
	_, err := db.Exec("INSERT INTO organizations (name, address) VALUES ($1, $2)", name, address)
	return err
}

func (db *DB) SearchEmployeesByQuery(query string) ([]models.Employee, error) {
	if query == "" {
		query = "%"
	} else {
		query = "%" + query + "%"
	}

	sqlQuery := `
        SELECT e.id, e.first_name, e.last_name, e.organization_id, o.name
        FROM employees e
        LEFT JOIN organizations o ON e.organization_id = o.id
        WHERE e.first_name ILIKE $1
           OR e.last_name ILIKE $1
           OR o.name ILIKE $1
    `
	rows, err := db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.OrganizationID, &emp.OrganizationName); err != nil {
			log.Println("Ошибка чтения строки:", err)
			continue
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (db *DB) DeleteEmployee(id string) error {
	_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
	return err
}
