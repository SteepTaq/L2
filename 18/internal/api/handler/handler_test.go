package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wb_l2/18/internal/repository"
	"wb_l2/18/internal/service"
)

func setupTestHandler() *Handler {
	repo := repository.NewRepository(repository.InMemory)
	svc := service.NewService(repo)
	h := NewHandler(svc)
	RegisterHandlers(h)
	return h
}

func TestPing(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Pong!" {
		t.Errorf("Expected message 'Pong!', got %v", response["message"])
	}
}

func TestNotFound(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Unknown endpoint" {
		t.Errorf("Expected error 'Unknown endpoint', got %v", response["error"])
	}
}

func TestCreateEvent_Success(t *testing.T) {
	handler := setupTestHandler()

	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Event created successfully" {
		t.Errorf("Expected message 'Event created successfully', got %v", response["message"])
	}

	// Check that data contains id
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field in response")
	}

	if _, exists := data["id"]; !exists {
		t.Error("Expected id in response data")
	}
}

func TestCreateEvent_InvalidContentType(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer([]byte("test")))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d, got %d", http.StatusUnsupportedMediaType, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Content-Type must be application/json" {
		t.Errorf("Expected error about Content-Type, got %v", response["error"])
	}
}

func TestCreateEvent_InvalidFormat(t *testing.T) {
	handler := setupTestHandler()

	// Missing required fields
	eventData := map[string]interface{}{
		"name": "Test Event",
		// missing date and user_id
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Invalid body format" {
		t.Errorf("Expected error 'Invalid body format', got %v", response["error"])
	}
}

func TestCreateEvent_WrongMethod(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/create_event", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestListEventsForDay_Success(t *testing.T) {
	handler := setupTestHandler()

	// First create an event
	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.mux.ServeHTTP(w, req)

	// Now list events for that day
	req = httptest.NewRequest("GET", "/events_for_day?user_id=1&date=2024-01-15", nil)
	w = httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "List of events for a day" {
		t.Errorf("Expected message 'List of events for a day', got %v", response["message"])
	}

	// Check that data contains events array
	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data field to be array in response")
	}

	if len(data) != 1 {
		t.Errorf("Expected 1 event, got %d", len(data))
	}
}

func TestListEventsForDay_InvalidQuery(t *testing.T) {
	handler := setupTestHandler()

	// Missing user_id
	req := httptest.NewRequest("GET", "/events_for_day?date=2024-01-15", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Invalid query provided" {
		t.Errorf("Expected error 'Invalid query provided', got %v", response["error"])
	}
}

func TestListEventsForDay_WrongMethod(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("POST", "/events_for_day", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestListEventsForWeek_Success(t *testing.T) {
	handler := setupTestHandler()

	// First create an event
	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.mux.ServeHTTP(w, req)

	// Now list events for that week
	req = httptest.NewRequest("GET", "/events_for_week?user_id=1&date=2024-01-15", nil)
	w = httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "List of events for a week" {
		t.Errorf("Expected message 'List of events for a week', got %v", response["message"])
	}
}

func TestListEventsForMonth_Success(t *testing.T) {
	handler := setupTestHandler()

	// First create an event
	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.mux.ServeHTTP(w, req)

	// Now list events for that month
	req = httptest.NewRequest("GET", "/events_for_month?user_id=1&date=2024-01-15", nil)
	w = httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "List of events for a month" {
		t.Errorf("Expected message 'List of events for a month', got %v", response["message"])
	}
}

func TestUpdateEvent_Success(t *testing.T) {
	handler := setupTestHandler()

	// First create an event
	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.mux.ServeHTTP(w, req)

	// Get the created event ID
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	data := createResponse["data"].(map[string]interface{})
	eventID := int(data["id"].(float64))

	// Now update the event
	updateData := map[string]interface{}{
		"id":   eventID,
		"name": "Updated Event",
		"date": "2024-01-16",
	}

	jsonData, _ = json.Marshal(updateData)
	req = httptest.NewRequest("POST", "/update_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Event updated successfully" {
		t.Errorf("Expected message 'Event updated successfully', got %v", response["message"])
	}
}

func TestUpdateEvent_NotFound(t *testing.T) {
	handler := setupTestHandler()

	updateData := map[string]interface{}{
		"id":   999, // Non-existent ID
		"name": "Updated Event",
		"date": "2024-01-16",
	}

	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest("POST", "/update_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Event is not found" {
		t.Errorf("Expected error 'Event is not found', got %v", response["error"])
	}
}

func TestUpdateEvent_InvalidFormat(t *testing.T) {
	handler := setupTestHandler()

	// Invalid JSON
	req := httptest.NewRequest("POST", "/update_event", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	handler := setupTestHandler()

	// First create an event
	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.mux.ServeHTTP(w, req)

	// Get the created event ID
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	data := createResponse["data"].(map[string]interface{})
	eventID := int(data["id"].(float64))

	// Now delete the event
	deleteData := map[string]interface{}{
		"id": eventID,
	}

	jsonData, _ = json.Marshal(deleteData)
	req = httptest.NewRequest("POST", "/delete_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Event deleted successfully" {
		t.Errorf("Expected message 'Event deleted successfully', got %v", response["message"])
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	handler := setupTestHandler()

	deleteData := map[string]interface{}{
		"id": 999, // Non-existent ID
	}

	jsonData, _ := json.Marshal(deleteData)
	req := httptest.NewRequest("POST", "/delete_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Event is not found" {
		t.Errorf("Expected error 'Event is not found', got %v", response["error"])
	}
}

func TestDeleteEvent_InvalidFormat(t *testing.T) {
	handler := setupTestHandler()

	// Invalid JSON
	req := httptest.NewRequest("POST", "/delete_event", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Additional edge case tests
func TestInvalidDate(t *testing.T) {
	handler := setupTestHandler()

	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "invalid-date",
		"user_id": 1,
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInvalidUserID(t *testing.T) {
	handler := setupTestHandler()

	eventData := map[string]interface{}{
		"name":    "Test Event",
		"date":    "2024-01-15",
		"user_id": 0, // Invalid user_id
	}

	jsonData, _ := json.Marshal(eventData)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListEvents_InvalidDate(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/events_for_day?user_id=1&date=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListEvents_InvalidUserID(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/events_for_day?user_id=invalid&date=2024-01-15", nil)
	w := httptest.NewRecorder()

	handler.mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
