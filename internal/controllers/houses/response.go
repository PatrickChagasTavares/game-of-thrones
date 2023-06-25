package houses

import (
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/houses"
)

func responseErr(err error, f func(int, any)) {
	switch err {
	case houses.ErrFind, houses.ErrNameUsed, houses.ErrHouseNotFound:
		f(http.StatusBadRequest, entities.NewHttpErr(http.StatusBadRequest, err.Error(), nil))
		return
	default:
		f(http.StatusInternalServerError, err)
	}
}