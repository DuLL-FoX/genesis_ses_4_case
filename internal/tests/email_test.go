package tests

import (
	"awesomeProject/internal/email"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	tests := []struct {
		to      string
		rate    float64
		wantErr error
	}{
		{"test@domain.com", 28.35, nil},
		{"", 28.35, errors.New("failed to send email: <error message>")},
	}

	for _, tt := range tests {
		err := email.SendEmail(tt.to, tt.rate)
		if tt.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
