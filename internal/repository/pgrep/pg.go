package pg

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type PgRepository struct {
	pgx *pgx.Conn
}

func (r *PgRepository) GetList() ([]string, error) {
	var keys []string

	rows, err := r.pgx.Query(context.TODO(), `
	SELECT id FROM metrics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *PgRepository) GetMetric(metricType string, metricName string) (res string, err error) {
	var metric models.Metrics

	err = r.pgx.QueryRow(context.TODO(), `
	SELECT id, m_type, delta, value
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": metricType, "metricName": metricName}).Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)

	if err != nil {
		return
	}

	switch metric.MType {
	case models.Counter:
		res = fmt.Sprint(*metric.Delta)
	case models.Gauge:
		res = fmt.Sprint(*metric.Value)
	default:
		err = models.ErrWrongMetric
	}

	return
}

func (r *PgRepository) GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error) {
	metricBody = models.Metrics{
		ID:    body.ID,
		MType: body.MType,
	}

	err = r.pgx.QueryRow(context.TODO(), `
	SELECT delta, value
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": body.MType, "metricName": body.ID}).Scan(&metricBody.Delta, &metricBody.Value)

	return
}

func (r *PgRepository) UpdateCounter(key string, delta int64) error {
	var id string
	err := r.pgx.QueryRow(context.TODO(), `
	SELECT id
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": models.Counter, "metricName": key}).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return r.insertCounter(key, delta)
		}
		return err
	}

	return r.updateCounter(key, delta)
}

func (r *PgRepository) UpdateGauge(key string, value float64) error {
	var id string
	err := r.pgx.QueryRow(context.TODO(), `
	SELECT id
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": models.Gauge, "metricName": key}).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return r.insertGauge(key, value)
		}
		return err
	}

	return r.updateGauge(key, value)
}

func (r *PgRepository) Close() {
	r.pgx.Close(context.Background())
}

func (r *PgRepository) Ping(ctx context.Context) error {
	return r.pgx.Ping(ctx)
}

func (r *PgRepository) insertCounter(name string, delta int64) error {
	_, err := r.pgx.Exec(context.TODO(), `
	INSERT INTO metrics (id, m_type, delta)
	VALUES (@metricName, @metricType, @delta)
	`, pgx.NamedArgs{
		"metricType": models.Counter,
		"metricName": name,
		"delta":      delta,
	})
	return err
}

func (r *PgRepository) insertGauge(name string, value float64) error {
	_, err := r.pgx.Exec(context.TODO(), `
	INSERT INTO metrics (id, m_type, value)
	VALUES (@metricName, @metricType, @value)
	`, pgx.NamedArgs{
		"metricType": models.Gauge,
		"metricName": name,
		"value":      value,
	})
	return err
}

func (r *PgRepository) updateCounter(name string, delta int64) error {
	_, err := r.pgx.Exec(context.TODO(), `
	UPDATE metrics
	SET delta = delta + @delta
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"delta": delta, "metricType": models.Counter, "metricName": name})

	return err
}

func (r *PgRepository) updateGauge(name string, value float64) error {
	_, err := r.pgx.Exec(context.TODO(), `
	UPDATE metrics
	SET value = @value
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"value": value, "metricType": models.Gauge, "metricName": name})

	return err
}

func NewPgRepository(pgConnect string) (*PgRepository, error) {

	connection, err := pgx.Connect(context.Background(), pgConnect)
	if err != nil {
		return nil, err
	}

	if err = migration(pgConnect); err != nil {
		return nil, err
	}
	rep := PgRepository{
		pgx: connection,
	}

	return &rep, nil
}

func migration(pgConnect string) error {
	db, err := sql.Open("postgres", pgConnect)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	return migration.Up()
}
