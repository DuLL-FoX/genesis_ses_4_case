package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"awesomeProject/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeHandler(t *testing.T) {
	form := strings.NewReader("email=test@domain.com")
	req, err := http.NewRequest("POST", "/subscribe", form)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.SubscribeHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Subscribed email: test@domain.com")
}
