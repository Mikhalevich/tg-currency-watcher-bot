package driver

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PGX struct {
}

func NewPGX() *PGX {
	return &PGX{}
}

func (p *PGX) IsConstraintError(err error, constraint string) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == constraint {
			return true
		}
	}

	return false
}
