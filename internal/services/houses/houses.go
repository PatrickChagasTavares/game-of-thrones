package houses

import (
	"context"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	IService interface {
		Create(ctx context.Context, newHouse entities.HouseRequest) (id string, err error)
		Find(ctx context.Context, name string) (houses []entities.House, err error)
		FindByID(ctx context.Context, id string) (house entities.House, err error)
		Update(ctx context.Context, updateHouse entities.HouseRequest) (house entities.House, err error)
		Delete(ctx context.Context, id string) (err error)
	}

	services struct {
		repositories *repositories.Container
		log          logger.Logger
	}
)

func New(repo *repositories.Container, log logger.Logger) IService {
	return &services{repositories: repo, log: log}
}

func (srv *services) Create(ctx context.Context, newHouse entities.HouseRequest) (id string, err error) {
	if _, err := srv.repositories.Database.House.FindByName(ctx, newHouse.Name); err == nil {
		return id, ErrNameUsed
	}

	newHouse.PreSave()

	err = srv.repositories.Database.House.Create(ctx, newHouse)
	if err != nil {
		srv.log.Error("Srv.Find: ", "create house ", err, ", playload: ", newHouse)
		return id, err
	}

	return newHouse.ID, nil
}

func (srv *services) Find(ctx context.Context, name string) (houses []entities.House, err error) {
	if len(name) > 0 {
		house, err := srv.repositories.Database.House.FindByName(ctx, name)
		if err != nil {
			srv.log.Error("Srv.Find: ", "House not found by name ", name)
			return nil, ErrFind
		}

		return []entities.House{house}, nil
	}

	houses, err = srv.repositories.Database.House.Find(ctx)
	if err != nil {
		srv.log.Error("Srv.Find: ", "Houses not found ", err)
		return nil, ErrFind
	}

	return houses, nil
}

func (srv *services) FindByID(ctx context.Context, id string) (house entities.House, err error) {
	house, err = srv.repositories.Database.House.FindByID(ctx, id)
	if err != nil {
		srv.log.Error("Srv.FindByID: ", "House not found ", id)
		return house, err
	}

	return house, nil
}

func (srv *services) Update(ctx context.Context, updateHouse entities.HouseRequest) (house entities.House, err error) {
	house, err = srv.FindByID(ctx, updateHouse.ID)
	if err != nil {
		return
	}

	if _, err := srv.repositories.Database.House.FindByName(ctx, updateHouse.Name); err == nil {
		return house, ErrNameUsed
	}

	house.PreUpdate(updateHouse)

	err = srv.repositories.Database.House.Update(ctx, &house)
	if err != nil {
		return house, err
	}

	return house, nil
}

func (srv *services) Delete(ctx context.Context, id string) (err error) {
	_, err = srv.repositories.Database.House.FindByID(ctx, id)
	if err != nil {
		return ErrHouseNotFound
	}

	err = srv.repositories.Database.House.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
