package main

import (
	"testing"
)

func expectSuccess(t *testing.T, args []string) {
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestNoArgs(t *testing.T) {
	expectSuccess(t, []string{})
}
