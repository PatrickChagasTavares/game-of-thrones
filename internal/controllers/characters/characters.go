package characters

import (
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
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
	ctx, span := tracer.Span(c.Context(), "controllers.characters.create")
	defer span.End()

	var newCharacter entities.CharacterRequest
	if err := c.Decode(&newCharacter); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrDecode)
		return
	}

	if err := c.Validate(newCharacter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err := ctrl.srv.Character.Create(ctx, newCharacter)
	if err != nil {
		ctrl.log.Error("Ctrl.Create: ", "Error on create character: ", newCharacter)
		responseErr(ctx, err, c.JSON)
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
	ctx, span := tracer.Span(c.Context(), "controllers.characters.find")
	defer span.End()

	characters, err := ctrl.srv.Character.Find(ctx)
	if err != nil {
		ctrl.log.Error("Ctrl.Find: ", "Error on find characters: ", err)
		responseErr(ctx, err, c.JSON)
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
	ctx, span := tracer.Span(c.Context(), "controllers.characters.findbyid")
	defer span.End()

	id := c.GetParam("id")

	characters, err := ctrl.srv.Character.FindByID(ctx, id)
	if err != nil {
		ctrl.log.Error("Ctrl.FindByID: ", "Error on find character: ", id)
		responseErr(ctx, err, c.JSON)
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
	ctx, span := tracer.Span(c.Context(), "controllers.characters.update")
	defer span.End()

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

	characters, err := ctrl.srv.Character.Update(ctx, updateCharacter)
	if err != nil {
		ctrl.log.Error("Ctrl.Update: ", "Error on update character: ", updateCharacter)
		responseErr(ctx, err, c.JSON)
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
	ctx, span := tracer.Span(c.Context(), "controllers.characters.delete")
	defer span.End()

	id := c.GetParam("id")

	err := ctrl.srv.Character.Delete(ctx, id)
	if err != nil {
		ctrl.log.Error("Ctrl.Delete: ", "Error on delete character: ", id)
		responseErr(ctx, err, c.JSON)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
