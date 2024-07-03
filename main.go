package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	defer closeDB()

	r := mux.NewRouter()
	r.HandleFunc("/todos", createTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id}", getTodoHandler).Methods("GET")
	r.HandleFunc("/todos/{id}", updateTodoHandler).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodoHandler).Methods("DELETE")
	r.HandleFunc("/todos", listTodosHandler).Methods("GET")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
