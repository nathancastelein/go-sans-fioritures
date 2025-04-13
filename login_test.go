package main

import (
	"testing"
)

func TestUser_LogValue(t *testing.T) {
	// Arrange
	user := &User{
		ID:        1,
		FirstName: "Peter",
		LastName:  "Parker",
		HeroName:  "Spider-Man",
	}
	expected := "[id=1 hero_name=Spider-Man]"

	// Act
	got := user.LogValue()

	// Assert
	if got.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, got.String())
	}
}
