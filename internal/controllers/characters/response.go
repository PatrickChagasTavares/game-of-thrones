package characters

import (
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/characters"
)

func responseErr(err error, f func(int, any)) {
	switch err {
	case characters.ErrFind, characters.ErrCharacterNotFound:
		f(http.StatusBadRequest, entities.NewHttpErr(http.StatusBadRequest, err.Error(), nil))
		return
	default:
		f(http.StatusInternalServerError, err)
	}
}
