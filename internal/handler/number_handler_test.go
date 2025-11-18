package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNumberService struct {
	mock.Mock
}

func (m *MockNumberService) AddNumber(value int) ([]int, error) {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func TestNumberHandler_AddNumber_Success(t *testing.T) {
	// Arrange
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	mockService.On("AddNumber", 3).Return([]int{1, 2, 3}, nil)

	reqBody := AddNumberRequest{Number: 3}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	handler.AddNumber(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response AddNumberResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, response.Numbers)
	mockService.AssertExpectations(t)
}

func TestNumberHandler_AddNumber_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	// Act
	handler.AddNumber(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "invalid request body")
}

func TestNumberHandler_AddNumber_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	mockService.On("AddNumber", 5).Return(nil, errors.New("database error"))

	reqBody := AddNumberRequest{Number: 5}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.AddNumber(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "failed to add number")
	mockService.AssertExpectations(t)
}

func TestNumberHandler_AddNumber_MethodNotAllowed(t *testing.T) {
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/numbers", nil)
	w := httptest.NewRecorder()

	handler.AddNumber(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "method not allowed")
}

func TestNumberHandler_HealthCheck_Success(t *testing.T) {
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestNumberHandler_HealthCheck_MethodNotAllowed(t *testing.T) {
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestNumberHandler_AddNumber_NegativeNumber(t *testing.T) {
	mockService := new(MockNumberService)
	handler := NewNumberHandler(mockService)

	mockService.On("AddNumber", -10).Return([]int{-10, 0, 5}, nil)

	reqBody := AddNumberRequest{Number: -10}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.AddNumber(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AddNumberResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, []int{-10, 0, 5}, response.Numbers)
	mockService.AssertExpectations(t)
}
