package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tasks/task-assessment/internal/api"
	"github.com/tasks/task-assessment/internal/db"
)

func main() {
	store := db.NewEmployeeStore()
	r := mux.NewRouter()
	apiRouter := api.NewAPI(store)
	r.PathPrefix("/").Handler((apiRouter))

	log.Println("Server is starting on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
