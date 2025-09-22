package pgerr

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGErrorClassification int

const (
	NonRetriable PGErrorClassification = iota
	Retriable
)

type PostgresErrorClassifier struct{}

func NewPostgresErrorClassifier() *PostgresErrorClassifier {
	return &PostgresErrorClassifier{}
}

func (c *PostgresErrorClassifier) Classify(err error) PGErrorClassification {
	if err == nil {
		return NonRetriable
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return classifyPgError(pgErr)
	}

	return NonRetriable
}

func classifyPgError(pgErr *pgconn.PgError) PGErrorClassification {

	switch pgErr.Code {
	// Класс 08 - Ошибки соединения
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection,
		pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
		pgerrcode.TransactionResolutionUnknown,
		pgerrcode.ProtocolViolation:
		return Retriable
	}

	return NonRetriable
}
