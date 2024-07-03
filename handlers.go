package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = gocql.TimeUUID().String()
	todo.Created = time.Now()
	todo.Updated = time.Now()

	if err := session.Query(`INSERT INTO todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var todo Todo
	if err := session.Query(`SELECT id, user_id, title, description, status, created, updated FROM todos WHERE id = ?`, id).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated); err != nil {
		if err == gocql.ErrNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Updated = time.Now()

	if err := session.Query(`UPDATE todos SET title = ?, description = ?, status = ?, updated = ? WHERE id = ?`,
		todo.Title, todo.Description, todo.Status, todo.Updated, id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := session.Query(`DELETE FROM todos WHERE id = ?`, id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	var todos []Todo
	query := `SELECT id, user_id, title, description, status, created, updated FROM todos WHERE user_id = ?`
	var iter *gocql.Iter

	if status != "" {
		query += " AND status = ? LIMIT ? OFFSET ? ALLOW FILTERING"
		iter = session.Query(query, userID, status, pageSizeInt, (pageInt-1)*pageSizeInt).Iter()
	} else {
		query += " LIMIT ? OFFSET ? ALLOW FILTERING"
		iter = session.Query(query, userID, pageSizeInt, (pageInt-1)*pageSizeInt).Iter()
	}

	for {
		var todo Todo
		if !iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
			break
		}
		todos = append(todos, todo)
	}

	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}
