package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	House struct {
		ID             string     `db:"id" json:"-"`
		Name           string     `db:"name" json:"name"`
		Region         string     `db:"region" json:"region"`
		FoundationYear string     `db:"foundation_year" json:"foundation_year"`
		CurrentLord    string     `db:"current_lord" json:"current_lord"`
		CreatedAt      *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
	}

	HouseRequest struct {
		ID             string `json:"-"`
		Name           string `json:"name" validate:"required,min=3,max=200"`
		Region         string `json:"region" validate:"required,min=3,max=100"`
		FoundationYear string `json:"foundation_year" validate:"required,min=1,max=5"`
		CurrentLord    string `json:"current_lord,omitempty"`
	}
)

func (hr *HouseRequest) PreSave() {
	hr.ID = uuid.NewString()
}
