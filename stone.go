package main

import (
	"context"
	"fmt"
	"slices"
)

type StoneRepository interface {
	ListStones(ctx context.Context) []*InfinityStone
	GetStone(ctx context.Context, stoneName string) (*InfinityStone, error)
}

type InfinityStone struct {
	Name   string `json:"name"`
	Color  string `json:"color"`
	Power  string `json:"power"`
	Status string `json:"status"` // e.g. "secured", "missing"
}

type InMemoryStoneRepository struct {
	stones []*InfinityStone
}

func NewInMemoryStoneRepository() StoneRepository {
	return &InMemoryStoneRepository{
		stones: []*InfinityStone{
			{"space", "blue", "Teleportation", "secured"},
			{"mind", "yellow", "Mind Control", "secured"},
			{"reality", "red", "Reality Warping", "missing"},
			{"power", "purple", "Unlimited Strength", "secured"},
			{"time", "green", "Time Travel", "secured"},
			{"soul", "orange", "Soul Manipulation", "unknown"},
		},
	}
}

func (i *InMemoryStoneRepository) ListStones(ctx context.Context) []*InfinityStone {
	return i.stones
}

func (i *InMemoryStoneRepository) GetStone(ctx context.Context, stoneName string) (*InfinityStone, error) {
	stoneIndex := slices.IndexFunc(i.stones, func(stone *InfinityStone) bool {
		return stone.Name == stoneName
	})

	if stoneIndex == -1 {
		return nil, fmt.Errorf("stone %s not found", stoneName)
	}

	return i.stones[stoneIndex], nil
}
