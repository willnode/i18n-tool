package main

import (
	"log"
	"os"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	err = RealMain(currentDir + "/test")

	if err != nil {
		t.Fatalf(`Error %+v`, err)
	}
}
