package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var (
	todos         []Todo
	nextTodoID    int
	isFlagEnabled int32
	mutex         sync.Mutex
)

func disableFlag() {
	atomic.StoreInt32(&isFlagEnabled, 0)
	fmt.Println("Flag disabled.")
}

func isEndpointEnabled() bool {
	return atomic.LoadInt32(&isFlagEnabled) == 1
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if !isEndpointEnabled() {
		http.Error(w, "Endpoint disabled", http.StatusForbidden)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	newTodo.ID = nextTodoID
	nextTodoID++
	newTodo.Status = "Incomplete"
	todos = append(todos, newTodo)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func getAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	if !isEndpointEnabled() {
		http.Error(w, "Endpoint disabled", http.StatusForbidden)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func markCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if !isEndpointEnabled() {
		http.Error(w, "Endpoint disabled", http.StatusForbidden)
		return
	}

	todoID := r.URL.Query().Get("id")
	if todoID == "" {
		http.Error(w, "Missing todo ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID parameter", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var updatedTodo Todo
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Status = "Complete"
			updatedTodo = todos[i]
			break
		}
	}

	if updatedTodo.ID == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if !isEndpointEnabled() {
		http.Error(w, "Endpoint disabled", http.StatusForbidden)
		return
	}

	todoID := r.URL.Query().Get("id")
	if todoID == "" {
		http.Error(w, "Missing todo ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID parameter", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var deletedTodo Todo
	var newTodos []Todo
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		} else {
			deletedTodo = todo
		}
	}

	if deletedTodo.ID == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos = newTodos

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedTodo)
}

func toggleFlagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	atomic.StoreInt32(&isFlagEnabled, 1-isFlagEnabled) // Toggle the flag

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/create", createTodoHandler)
	http.HandleFunc("/todos", getAllTodosHandler)
	http.HandleFunc("/complete", markCompleteHandler)
	http.HandleFunc("/delete", deleteTodoHandler)
	http.HandleFunc("/toggleflag", toggleFlagHandler)

	disableFlag() // Enable the flag initially

	fmt.Println("Server is running on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic(err)
	}
}
