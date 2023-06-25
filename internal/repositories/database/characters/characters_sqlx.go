package characters

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

func (repo *repoSqlx) Create(ctx context.Context, character entities.CharacterRequest) (err error) {
	_, err = repo.writer.ExecContext(ctx,
		`INSERT INTO characters 
		(id,name,tv_series,created_at)
		VALUES ($1, $2, $3, $4);`,
		character.ID, character.Name, character.TVSeries, character.CreatedAt)
	if err != nil {
		repo.log.ErrorContext(ctx, "characters.SqlxRepo.Create", err)
		return errors.New("problem to create character")
	}

	return nil
}

func (repo *repoSqlx) Find(ctx context.Context) (characters []entities.Character, err error) {
	characters = make([]entities.Character, 0)
	query := `
	SELECT id, name, tv_series, created_at, updated_at
	FROM characters
	WHERE deleted_at is null
	ORDER BY created_at DESC;
	`
	err = repo.reader.SelectContext(ctx, &characters, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return characters, nil
		}
		repo.log.ErrorContext(ctx, "characters.SqlxRepo.Find", "Error on find characters: ", err)
		return nil, errors.New("problem to find characters")
	}

	return characters, nil
}

func (repo *repoSqlx) FindByID(ctx context.Context, id string) (character entities.Character, err error) {
	query := `
	SELECT id, name, tv_series, created_at, updated_at
	FROM characters
	WHERE id =$1 AND deleted_at is null;`
	err = repo.reader.GetContext(ctx, &character, query, id)
	if err != nil {
		repo.log.ErrorContext(ctx, "characters.SqlxRepo.FindByID", "Error on find character by id: ", id, err)
		return character, errors.New("character is not found or deleted")
	}

	return character, nil
}

func (repo *repoSqlx) Update(ctx context.Context, character *entities.Character) (err error) {
	query := `
	UPDATE characters
	SET name = :name, tv_series = :tv_series, updated_at = :updated_at
	WHERE id = :id;
	`
	_, err = repo.writer.NamedExecContext(ctx, query, character)
	if err != nil {
		repo.log.ErrorContext(ctx, "characters.SqlxRepo.Update", "Error on update character: ", character, err)
		return errors.New("failed to update character")
	}

	return nil
}

var timeNow = time.Now

func (repo *repoSqlx) Delete(ctx context.Context, id string) (err error) {
	query := `
	UPDATE characters
	SET deleted_at = $1
	WHERE id = $2;
	`
	_, err = repo.writer.ExecContext(ctx, query, timeNow(), id)
	if err != nil {
		repo.log.ErrorContext(ctx, "characters.SqlxRepo.Delete", "Error on delete character: ", id, err)
		return errors.New("failed to delete character")
	}

	return nil
}
