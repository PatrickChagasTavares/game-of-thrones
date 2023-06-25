package houses

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type repoSqlx struct {
	log    logger.Logger
	writer *sqlx.DB
	reader *sqlx.DB
}

func NewSqlx(log logger.Logger, writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{log: log, writer: writer, reader: reader}
}

func (repo *repoSqlx) Create(ctx context.Context, house entities.HouseRequest) (err error) {
	_, err = repo.writer.ExecContext(ctx,
		`INSERT INTO houses 
		(id,name,region,foundation_year,current_lord,created_at)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		house.ID, house.Name, house.Region, house.FoundationYear, house.CurrentLord, house.CreatedAt)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.Create", err)
		return errors.New("problem to create house")
	}

	return nil
}

func (repo *repoSqlx) Find(ctx context.Context) (houses []entities.House, err error) {
	houses = make([]entities.House, 0)
	query := `
	SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
	FROM houses
	WHERE deleted_at is null
	ORDER BY created_at DESC;
	`
	err = repo.reader.SelectContext(ctx, &houses, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return houses, nil
		}
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.FindByID", "Error on find house: ", err)
		return nil, errors.New("problem to find houses")
	}

	return houses, nil
}

func (repo *repoSqlx) FindByID(ctx context.Context, id string) (houses entities.House, err error) {
	query := `
	SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
	FROM houses
	WHERE id =$1 AND deleted_at is null;`
	err = repo.reader.GetContext(ctx, &houses, query, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.FindByID", "Error on find house by id: ", id, err)
		return houses, errors.New("house is not found or deleted")
	}

	return houses, nil
}

func (repo *repoSqlx) FindByName(ctx context.Context, name string) (houses entities.House, err error) {
	query := `
	SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
	FROM houses
	WHERE name=$1 AND deleted_at is null;`
	err = repo.reader.GetContext(ctx, &houses, query, name)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.FindByName", "Error on find house by name: ", name, err)
		return houses, errors.New("house is not found or deleted")
	}

	return houses, nil
}

func (repo *repoSqlx) Update(ctx context.Context, house *entities.House) (err error) {
	query := `
	UPDATE houses
	SET name = :name, region = :region, foundation_year = :foundation_year, current_lord = :current_lord, updated_at = :updated_at
	WHERE id = :id;
	`
	_, err = repo.writer.NamedExecContext(ctx, query, house)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.Update", "Error on update house: ", house, err)
		return errors.New("failed to delete house")
	}

	return nil
}

func (repo *repoSqlx) Delete(ctx context.Context, id string) (err error) {
	query := `
	UPDATE houses
	SET deleted_at = $1
	WHERE id = $2;
	`
	now := time.Now()
	_, err = repo.writer.ExecContext(ctx, query, now, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.Delete", "Error on delete house: ", id, err)
		return errors.New("failed to delete house")
	}

	return nil
}
