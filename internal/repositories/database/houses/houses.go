//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go
package houses

import (
	"context"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
)

type IRepository interface {
	Create(ctx context.Context, house entities.HouseRequest) (err error)
	Find(ctx context.Context) (houses []entities.House, err error)
	FindByID(ctx context.Context, id string) (houses entities.House, err error)
	FindByName(ctx context.Context, name string) (houses entities.House, err error)
	RemoveLord(ctx context.Context, lordID string) (err error)
	Update(ctx context.Context, house *entities.House) (err error)
	Delete(ctx context.Context, id string) (err error)
}
