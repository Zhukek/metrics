package pg

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
)

type PgRepository struct {
	DB *sql.DB
}

/* func (r *PgRepository) GetList() []string {

}

func (r *PgRepository) GetMetric(metricType string, metricName string) (res string, err error) {

}

func (r *PgRepository) GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error) {

}

func (r *PgRepository) UpdateCounter(key string, value int64) {

}

func (r *PgRepository) UpdateGauge(key string, value float64) {

}

func (r *PgRepository) GetAllMetrics() map[string]models.Metrics {

} */

func (r *PgRepository) Close() {
	r.DB.Close()
}

func NewPgRepository(pgConnect string) ( /* repository.Repository */ *PgRepository, error) {

	db, err := sql.Open("pgx", pgConnect)
	if err != nil {
		return nil, err
	}

	dbdriver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "pgx", dbdriver)
	if err != nil {
		return nil, err
	}
	err = migration.Up()
	if err != nil {
		return nil, err
	}
	rep := PgRepository{
		DB: db,
	}

	return &rep, nil
}
