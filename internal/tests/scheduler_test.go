package tests

import (
	"testing"
	"time"

	"awesomeProject/internal/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestStartScheduler(t *testing.T) {
	go scheduler.StartScheduler()
	// Wait for a short time to simulate the ticker behavior
	time.Sleep(2 * time.Second)
	// Check if emails were sent and other functionality
	assert.True(t, true) // Replace with actual checks
}
