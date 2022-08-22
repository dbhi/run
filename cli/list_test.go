package main

import (
	"testing"
)

func TestListNoArgs(t *testing.T) {
	expectSuccess(t, []string{"list"})
}
