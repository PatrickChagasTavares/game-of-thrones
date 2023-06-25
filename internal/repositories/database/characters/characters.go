//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go
package characters

import (
	"context"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
)

type IRepository interface {
	Create(ctx context.Context, character entities.CharacterRequest) (err error)
	Find(ctx context.Context) (characters []entities.Character, err error)
	FindByID(ctx context.Context, id string) (characters entities.Character, err error)
	Update(ctx context.Context, character *entities.Character) (err error)
	Delete(ctx context.Context, id string) (err error)
}
