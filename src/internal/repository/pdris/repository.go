package pdris

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"gihub.com/bongerka/sberPDRIS/internal/model"
	def "gihub.com/bongerka/sberPDRIS/internal/repository"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) def.PdrisRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdateValue(ctx context.Context, value int) error {
	sql := `INSERT INTO pdris.value (key, value) values ($1, $2)
		    ON CONFLICT (key) DO UPDATE 
		    SET value = excluded.value`

	_, err := r.db.Exec(ctx, sql, model.POSTGRES_KEY, value)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetValue(ctx context.Context) (int, error) {
	sql := `SELECT value
            FROM pdris.value
            WHERE key = $1`

	rows, err := r.db.Query(ctx, sql, model.POSTGRES_KEY)
	if err != nil {
		return 0, err
	}

	value, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int, error) {
		var value int
		if err = row.Scan(&value); err != nil {
			return 0, err
		}

		return value, nil
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return value, nil
}
