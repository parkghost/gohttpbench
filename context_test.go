package main

import (
	"testing"
)

func TestSetAndGetString(t *testing.T) {
	context := NewContext(&Config{})
	key := "key"
	value := "value"
	context.SetString(key, value)
	if context.GetString(key) != value {
		t.Fatalf("expected %s, got %s", key, value)
	}
}

func TestSetAndGetInt(t *testing.T) {
	context := NewContext(&Config{})
	key := "key"
	value := 123
	context.SetInt(key, value)
	if context.GetInt(key) != value {
		t.Fatalf("expected %s, got %s", key, value)
	}
}
