package houses

import (
	"context"
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/houses"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
)

func responseErr(ctx context.Context, err error, f func(int, any)) {
	_, span := tracer.Span(ctx, "controllers.houses.responseErr")
	defer span.End()

	switch err {
	case houses.ErrFind, houses.ErrNameUsed, houses.ErrHouseNotFound:
		f(http.StatusBadRequest, entities.NewHttpErr(http.StatusBadRequest, err.Error(), nil))
		return
	default:
		f(http.StatusInternalServerError, err.Error())
	}
}
