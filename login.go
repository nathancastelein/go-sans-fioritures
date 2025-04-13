package main

import (
	"context"
	"errors"
	"log/slog"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	HeroName  string
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", u.ID),
		slog.String("hero_name", u.HeroName),
	)
}

type Login interface {
	LogUser(ctx context.Context, username string, password string) (*User, error)
}

type InMemoryLogin struct{}

func NewInMemoryLogin() Login {
	return &InMemoryLogin{}
}

func (i *InMemoryLogin) LogUser(ctx context.Context, username string, password string) (*User, error) {
	if username != "tony.stark" || password != "howard" {
		return nil, errors.New("invalid login/password")
	}
	return &User{
		ID:        1,
		FirstName: "Tony",
		LastName:  "Stark",
		HeroName:  "Iron Man",
	}, nil
}
