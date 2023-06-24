package houses

import (
	"context"
	"database/sql"
	"errors"

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
		(id,name,region,foundation_year,current_lord)
		VALUES ($1, $2, $3, $4, $5);`,
		house.ID, house.Name, house.Region, house.FoundationYear, house.CurrentLord)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.Create", err)
		return errors.New("problem to create house")
	}

	return nil
}

func (repo *repoSqlx) Find(ctx context.Context, limit, offset uint) (houses []entities.House, err error) {
	query := `
	SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
	FROM houses
	WHERE deleted_at is null
	ORDER BY created_at DESC
	LIMIT $1
	OFFSET $2;
	`
	err = repo.reader.SelectContext(ctx, &houses, query, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entities.House{}, nil
		}
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.FindByID", "Error on find house: limit - ", limit, ", offset -", offset, err)
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

func (repo *repoSqlx) Delete(ctx context.Context, id string) (err error) {
	query := `
	UPDATE houses
	SET deleted_at = CURRENT_TIMESTAMP AND updated_at = CURRENT_TIMESTAMP
	WHERE id = $1;
	`
	_, err = repo.writer.ExecContext(ctx, query, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "houses.SqlxRepo.Delete", "Error on delete house: ", id, err)
		return errors.New("failed to delete")
	}

	return nil
}
