package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/hablof/generate-random-value/internal/models"
)

const (
	tableVals    = "vals"
	colID        = "id"
	colVal       = "val"
	colRequestID = "request_id"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository struct {
	db            *sql.DB
	initStatement squirrel.StatementBuilderType
}

func NewRepository(db *sql.DB) *Repository {
	initStatement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	return &Repository{
		db:            db,
		initStatement: initStatement,
	}
}

func (r *Repository) Create(unit models.RandomValue) (uint64, error) {
	query := r.initStatement.Insert(tableVals).Columns(colVal)
	if unit.RequestID.IsValid {
		query = query.Columns(colRequestID).Values(unit.Value, unit.RequestID.S)
	} else {
		query = query.Values(unit.Value)
	}

	queryString, args, err := query.Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	var id uint64
	row := r.db.QueryRowContext(ctx, queryString, args...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) ReadByReqID(reqID string) (models.RandomValue, error) {
	queryString, args, err := r.initStatement.
		Select(colID, colVal).
		From(tableVals).
		Where(squirrel.Eq{colRequestID: reqID}).
		ToSql()

	if err != nil {
		return models.RandomValue{}, err
	}

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	unit := models.RandomValue{}
	row := r.db.QueryRowContext(ctx, queryString, args...)
	switch err := row.Scan(&unit.ID, &unit.Value); {
	case errors.Is(err, sql.ErrNoRows):
		return models.RandomValue{}, ErrNotFound

	case err != nil:
		return models.RandomValue{}, err
	}

	return unit, nil
}

func (r *Repository) ReadByValID(valID int) (string, error) {
	queryString, args, err := r.initStatement.
		Select(colVal).
		From(tableVals).
		Where(squirrel.Eq{colID: valID}).
		ToSql()

	if err != nil {
		return "", err
	}

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	var value string
	row := r.db.QueryRowContext(ctx, queryString, args...)
	switch err := row.Scan(&value); {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrNotFound

	case err != nil:
		return "", err
	}

	return value, nil
}
