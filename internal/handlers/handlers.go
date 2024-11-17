package handlers

import (
	_ "database/sql"
	"employeeOrgDB/internal/db"
	"employeeOrgDB/internal/models"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var database *db.DB

func InitHandlers(databaseConn *db.DB) {
	database = databaseConn
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	employees, err := database.SearchEmployeesByQuery(query)
	if err != nil {
		log.Println("Ошибка выполнения запроса:", err)
		http.Error(w, "Ошибка получения данных сотрудников", http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено сотрудников: %d\n", len(employees))

	data := struct {
		Employees   []models.Employee
		SearchQuery string
	}{
		Employees:   employees,
		SearchQuery: query,
	}

	tmpl, err := template.ParseFiles("templates/employees.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func AddEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		organizationID, _ := strconv.Atoi(r.FormValue("organization_id"))

		err := database.AddEmployee(firstName, lastName, organizationID)
		if err != nil {
			http.Error(w, "Ошибка добавления сотрудника", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/employees", http.StatusSeeOther)
		return
	}

	organizations, err := database.GetOrganizations()
	if err != nil {
		http.Error(w, "Ошибка получения данных организаций", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/add_employee.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, organizations)
}

func OrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	organizations, err := database.GetOrganizations()
	if err != nil {
		http.Error(w, "Ошибка получения данных организаций", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/organizations.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, organizations)
}

func AddOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		address := r.FormValue("address")
		returnToAddEmployee := r.FormValue("return_to_add_employee")

		err := database.AddOrganization(name, address)
		if err != nil {
			http.Error(w, "Ошибка добавления организации", http.StatusInternalServerError)
			return
		}

		if returnToAddEmployee == "true" {
			redirectURL := "/employees/add?return_to_add_employee=false"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/organizations", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/add_organization.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	parseUrl, _ := url.Parse(r.URL.String())
	returnToAddEmployee := parseUrl.Query().Get("return_to_add_employee")

	data := struct {
		ReturnToAddEmployee string
	}{
		ReturnToAddEmployee: returnToAddEmployee,
	}

	tmpl.Execute(w, data)
}

func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	empID := r.URL.Query().Get("id")

	err := database.DeleteEmployee(empID)
	if err != nil {
		http.Error(w, "Ошибка удаления сотрудника", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/employees", http.StatusSeeOther)
}
