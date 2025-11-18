package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNumberRepository struct {
	mock.Mock
}

func (m *MockNumberRepository) Save(value int) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockNumberRepository) GetAllSorted() ([]int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func TestNumberService_AddNumber_Success(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	service := NewNumberService(mockRepo)

	mockRepo.On("Save", 5).Return(nil)
	mockRepo.On("GetAllSorted").Return([]int{1, 3, 5}, nil)

	result, err := service.AddNumber(5)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, result)
	mockRepo.AssertExpectations(t)
}

func TestNumberService_AddNumber_SaveError(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	service := NewNumberService(mockRepo)

	expectedError := errors.New("database error")
	mockRepo.On("Save", 5).Return(expectedError)

	result, err := service.AddNumber(5)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to save number")
	mockRepo.AssertExpectations(t)
}

func TestNumberService_AddNumber_GetAllSortedError(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	service := NewNumberService(mockRepo)

	mockRepo.On("Save", 5).Return(nil)
	expectedError := errors.New("query error")
	mockRepo.On("GetAllSorted").Return(nil, expectedError)

	result, err := service.AddNumber(5)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get sorted numbers")
	mockRepo.AssertExpectations(t)
}

func TestNumberService_AddNumber_EmptyDatabase(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	service := NewNumberService(mockRepo)

	mockRepo.On("Save", 1).Return(nil)
	mockRepo.On("GetAllSorted").Return([]int{1}, nil)

	result, err := service.AddNumber(1)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, result)
	mockRepo.AssertExpectations(t)
}

func TestNumberService_AddNumber_NegativeNumber(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	service := NewNumberService(mockRepo)

	mockRepo.On("Save", -5).Return(nil)
	mockRepo.On("GetAllSorted").Return([]int{-5, 1, 3}, nil)

	result, err := service.AddNumber(-5)

	assert.NoError(t, err)
	assert.Equal(t, []int{-5, 1, 3}, result)
	mockRepo.AssertExpectations(t)
}
