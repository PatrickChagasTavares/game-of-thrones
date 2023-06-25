package characters

import (
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	IController interface {
		Create(c httpRouter.Context)
		Find(c httpRouter.Context)
		FindByID(c httpRouter.Context)
		Update(c httpRouter.Context)
		Delete(c httpRouter.Context)
	}
	controllers struct {
		srv *services.Container
		log logger.Logger
	}
)

func New(srv *services.Container, log logger.Logger) IController {
	return &controllers{srv: srv, log: log}
}

// character swagger document
// @Description Create one character
// @Tags character
// @Accept json
// @Produce json
// @Param character body entities.CharacterRequest true "create new character"
// @Success 201
// @Failure 400 {object} entities.HttpErr
// @Failure 409 {object} entities.HttpErr
// @Failure 500
// @Security ApiKeyAuth
// @Router /characters [post]
func (ctrl *controllers) Create(c httpRouter.Context) {
	var newCharacter entities.CharacterRequest
	if err := c.Decode(&newCharacter); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrDecode)
		return
	}

	if err := c.Validate(newCharacter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err := ctrl.srv.Character.Create(c.Context(), newCharacter)
	if err != nil {
		ctrl.log.Error("Ctrl.Create: ", "Error on create character: ", newCharacter)
		responseErr(err, c.JSON)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"id": id,
	})
}

// character swagger document
// @Description Find characters
// @Tags character
// @Accept json
// @Produce json
// @Success 200 {object} []entities.Character
// @Failure 500
// @Security ApiKeyAuth
// @Router /characters [get]
func (ctrl *controllers) Find(c httpRouter.Context) {

	characters, err := ctrl.srv.Character.Find(c.Context())
	if err != nil {
		ctrl.log.Error("Ctrl.Find: ", "Error on find characters: ", err)
		responseErr(err, c.JSON)
		return
	}

	c.JSON(http.StatusOK, characters)
}

// character swagger document
// @Description find character by id
// @Tags character
// @Accept json
// @Produce json
// @Param id path string true "Character ID"
// @Success 200 {object} entities.Character
// @Failure 500
// @Security ApiKeyAuth
// @Router /characters/:id [get]
func (ctrl *controllers) FindByID(c httpRouter.Context) {

	id := c.GetParam("id")

	characters, err := ctrl.srv.Character.FindByID(c.Context(), id)
	if err != nil {
		ctrl.log.Error("Ctrl.FindByID: ", "Error on find character: ", id)
		responseErr(err, c.JSON)
		return
	}

	c.JSON(http.StatusOK, characters)
}

// character swagger document
// @Description Update character
// @Tags character
// @Accept json
// @Produce json
// @Param id path string true "Character ID"
// @Param character body entities.CharacterRequest true "create new character"
// @Success 200 {object} entities.Character
// @Failure 400 {object} entities.HttpErr
// @Failure 500
// @Security ApiKeyAuth
// @Router /characters/:id [put]
func (ctrl *controllers) Update(c httpRouter.Context) {

	var updateCharacter entities.CharacterRequest
	if err := c.Decode(&updateCharacter); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrDecode)
		return
	}

	if err := c.Validate(updateCharacter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	updateCharacter.ID = c.GetParam("id")

	characters, err := ctrl.srv.Character.Update(c.Context(), updateCharacter)
	if err != nil {
		ctrl.log.Error("Ctrl.Update: ", "Error on update character: ", updateCharacter)
		responseErr(err, c.JSON)
		return
	}

	c.JSON(http.StatusOK, characters)
}

// character swagger document
// @Description Delete character
// @Tags character
// @Accept json
// @Produce json
// @Param id path string true "Character ID"
// @Param character body entities.CharacterRequest true "create new character"
// @Success 200 {object} entities.Character
// @Failure 400 {object} entities.HttpErr
// @Failure 500
// @Security ApiKeyAuth
// @Router /characters/:id [delete]
func (ctrl *controllers) Delete(c httpRouter.Context) {
	id := c.GetParam("id")
	if len(id) < 20 {
		c.JSON(http.StatusBadRequest, entities.NewHttpErr(http.StatusBadRequest, "id informad is invalid", id))
		return
	}

	err := ctrl.srv.Character.Delete(c.Context(), id)
	if err != nil {
		ctrl.log.Error("Ctrl.Delete: ", "Error on delete character: ", id)
		responseErr(err, c.JSON)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
