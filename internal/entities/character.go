package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type (
	Character struct {
		ID        string         `db:"id" json:"id"`
		Name      string         `db:"name" json:"name"`
		TVSeries  pq.StringArray `db:"tv_series" json:"tv_series"`
		CreatedAt time.Time      `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time     `db:"updated_at" json:"updated_at"`
	}

	CharacterRequest struct {
		ID        string         `json:"id"`
		Name      string         `json:"name" validate:"required,min=3,max=200"`
		TVSeries  pq.StringArray `json:"tv_series" validate:"required,min=1"`
		CreatedAt time.Time      `json:"-"`
	}
)

func (lr *CharacterRequest) PreSave() {
	lr.ID = uuid.NewString()
	lr.CreatedAt = time.Now()
}

func (l *Character) PreUpdate(character CharacterRequest) {
	if character.Name != l.Name {
		l.Name = character.Name
	}

	actSession := strings.Join(l.TVSeries, ",")
	newSession := strings.Join(character.TVSeries, ",")
	if actSession != newSession {
		l.TVSeries = character.TVSeries
	}

	now := time.Now()
	l.UpdatedAt = &now
}
