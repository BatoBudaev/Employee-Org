package handlers

import (
	_ "database/sql"
	"employeeOrgDB/internal/db"
	"html/template"
	"net/http"
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
	employees, err := database.GetEmployess()
	if err != nil {
		http.Error(w, "Ошибка получения данных сотрудников", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/employees.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, employees)
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
