package main

import (
	db2 "employeeOrgDB/internal/db"
	"employeeOrgDB/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := db2.InitDB("postgres", "1", "employee_org_database")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	log.Println("Сервер запущен на порту 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
