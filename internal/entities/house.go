package entities

import (
	"context"
	"time"

	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
	"github.com/google/uuid"
)

type (
	House struct {
		ID             string     `db:"id" json:"id"`
		Name           string     `db:"name" json:"name"`
		Region         string     `db:"region" json:"region"`
		FoundationYear string     `db:"foundation_year" json:"foundation_year"`
		CurrentLord    string     `db:"current_lord" json:"current_lord"`
		CreatedAt      time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
	}

	HouseRequest struct {
		ID             string     `json:"-"`
		Name           string     `json:"name" validate:"required,min=3,max=200"`
		Region         string     `json:"region" validate:"required,min=3,max=100"`
		FoundationYear string     `json:"foundation_year" validate:"required,min=1,max=5"`
		CurrentLord    string     `json:"current_lord,omitempty"`
		CreatedAt      time.Time  `db:"created_at" json:"-"`
		UpdatedAt      *time.Time `db:"updated_at" json:"-"`
	}
)

func (hr *HouseRequest) PreSave(ctx context.Context) {
	_, span := tracer.Span(ctx, "entities.house.presave")
	defer span.End()

	hr.ID = uuid.NewString()
	hr.CreatedAt = time.Now()
}

func (h *House) PreUpdate(ctx context.Context, house HouseRequest) {
	_, span := tracer.Span(ctx, "entities.house.preupdate")
	defer span.End()

	if house.Name != h.Name {
		h.Name = house.Name
	}

	if house.Region != h.Region {
		h.Region = house.Region
	}

	if house.FoundationYear != h.FoundationYear {
		h.FoundationYear = house.FoundationYear
	}

	if house.CurrentLord != h.CurrentLord {
		h.CurrentLord = house.CurrentLord
	}

	now := time.Now()
	h.UpdatedAt = &now
}
