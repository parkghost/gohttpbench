package main

import (
	"testing"
)

func TestSetAndGetString(t *testing.T) {
	key := "key"
	value := "value"

	context := NewContext(&Config{})
	context.SetString(key, value)

	got := context.GetString(key)
	if value != got {
		t.Fatalf("expected %s, got %s", value, got)
	}
}

func TestSetAndGetInt(t *testing.T) {
	key := "key"
	value := 123

	context := NewContext(&Config{})
	context.SetInt(key, value)

	got := context.GetInt(key)
	if value != got {
		t.Fatalf("expected %d, got %d", value, got)
	}
}
