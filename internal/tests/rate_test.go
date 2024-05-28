package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetExchangeRate is a mock function for GetExchangeRate
type MockGetExchangeRate struct {
	mock.Mock
}

func (m *MockGetExchangeRate) GetExchangeRate() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func TestGetExchangeRate(t *testing.T) {
	// Use the real function for this test
	rate, err := handlers.GetExchangeRate()
	assert.NoError(t, err)
	assert.Greater(t, rate, 0.0)
}

func TestExchangeRateHandler(t *testing.T) {
	// Create a mock for the GetExchangeRate function
	mockGetExchangeRate := new(MockGetExchangeRate)
	handlers.GetExchangeRate = mockGetExchangeRate.GetExchangeRate
	mockGetExchangeRate.On("GetExchangeRate").Return(28.35, nil)

	req, err := http.NewRequest("GET", "/rate", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExchangeRateHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Current USD to UAH exchange rate: 28.35")

	mockGetExchangeRate.AssertExpectations(t)
}
