package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCreateTodoHandler(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1) // Enable the endpoint
	ts := httptest.NewServer(http.HandlerFunc(createTodoHandler))
	defer ts.Close()
	client := ts.Client()

	payload := []byte(`{"title": "Test Todo"}`)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	var todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		t.Fatal(err)
	}

	if todo.Title != "Test Todo" || todo.Status != "Incomplete" {
		t.Errorf("Expected todo with title 'Test Todo' and status 'Incomplete', but got %+v", todo)
	}
}

func TestGetAllTodosHandler(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1) // Enable the endpoint
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(getAllTodosHandler))
	defer ts.Close()
	client := ts.Client()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var responseTodos []Todo
	err = json.NewDecoder(resp.Body).Decode(&responseTodos)
	if err != nil {
		t.Fatal(err)
	}

	if len(responseTodos) != 1 || responseTodos[0].Title != "Test Todo" {
		t.Errorf("Unexpected response todos: %+v", responseTodos)
	}
}

func TestMarkCompleteHandler(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1) // Enable the endpoint
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(markCompleteHandler))
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("PUT", ts.URL+"?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var updatedTodo Todo
	err = json.NewDecoder(resp.Body).Decode(&updatedTodo)
	if err != nil {
		t.Fatal(err)
	}

	if updatedTodo.Status != "Complete" {
		t.Errorf("Expected todo status 'Complete', but got '%s'", updatedTodo.Status)
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1) // Enable the endpoint
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(deleteTodoHandler))
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("DELETE", ts.URL+"?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var deletedTodo Todo
	err = json.NewDecoder(resp.Body).Decode(&deletedTodo)
	if err != nil {
		t.Fatal(err)
	}

	if deletedTodo.Title != "Test Todo" {
		t.Errorf("Expected deleted todo title 'Test Todo', but got '%s'", deletedTodo.Title)
	}

	if len(todos) != 0 {
		t.Errorf("Expected todos slice to be empty after deletion, but it's not")
	}
}

func TestToggleFlagHandler(t *testing.T) {
	resetState()
	ts := httptest.NewServer(http.HandlerFunc(toggleFlagHandler))
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("POST", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	if !isEndpointEnabled() {
		t.Error("Expected flag to be enabled after toggle, but it's not")
	}
}

func TestCreateTodoHandler_InvalidJSON(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1)
	ts := httptest.NewServer(http.HandlerFunc(createTodoHandler))
	defer ts.Close()
	client := ts.Client()

	payload := []byte(`invalid_json`)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetAllTodosHandler_FlagDisabled(t *testing.T) {
	resetState()
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(getAllTodosHandler))
	defer ts.Close()
	client := ts.Client()

	atomic.StoreInt32(&isFlagEnabled, 0) // Disable the endpoint

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, resp.StatusCode)
	}
}

func TestMarkCompleteHandler_MissingID(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1)
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(markCompleteHandler))
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("PUT", ts.URL, nil) // Missing ID parameter
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestDeleteTodoHandler_InvalidID(t *testing.T) {
	resetState()
	atomic.StoreInt32(&isFlagEnabled, 1)
	todos = append(todos, Todo{ID: 1, Title: "Test Todo", Status: "Incomplete"})
	ts := httptest.NewServer(http.HandlerFunc(deleteTodoHandler))
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("DELETE", ts.URL+"?id=2", nil) // Invalid ID
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestMainFunction(t *testing.T) {
	go main()
}

func resetState() {
	todos = nil
	nextTodoID = 1
	isFlagEnabled = 0
	mutex = sync.Mutex{}
}
