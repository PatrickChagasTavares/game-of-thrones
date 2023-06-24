package houses

import (
	"net/http"
	"strconv"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	IController interface {
		Create(c httpRouter.Context)
	}
	controllers struct {
		srv *services.Container
		log logger.Logger
	}
)

func New(srv *services.Container, log logger.Logger) IController {
	return &controllers{srv: srv, log: log}
}

// house swagger document
// @Description Create one house
// @Tags house
// @Param house body entities.HouseRequest true "create new house"
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {object} entities.HttpErr
// @Failure 409 {object} entities.HttpErr
// @Failure 500
// @Security ApiKeyAuth
// @Router /houses [post]
func (ctrl *controllers) Create(c httpRouter.Context) {
	var newHouse entities.HouseRequest
	if err := c.Decode(&newHouse); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrDecode)
		return
	}

	if err := c.Validate(newHouse); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err := ctrl.srv.House.Create(c.Context(), newHouse)
	if err != nil {
		ctrl.log.Error("Ctrl.Create: ", "Error on create house: ", newHouse)
		c.JSONError(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"id": id,
	})
}

// user swagger document
// @Description Create one house
// @Tags house
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Security ApiKeyAuth
// @Router /houses [get]
func (ctrl *controllers) Find(c httpRouter.Context) {
	limit, _ := strconv.ParseUint(c.GetQuery("limit"), 10, 64)
	offset, _ := strconv.ParseUint(c.GetQuery("offset"), 10, 64)

	houses, err := ctrl.srv.House.Find(c.Context(), uint(limit), uint(offset))
	if err != nil {
		ctrl.log.Error("Ctrl.Find: ", "Error on find houses: ", limit, offset)
		c.JSONError(err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"data":   houses,
		"limit":  limit,
		"offset": offset,
	})
}
