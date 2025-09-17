package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgRepository struct {
	pgx *pgx.Conn
}

type conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
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

func (r *PgRepository) GetMetric(metricType models.MType, metricName string) (res string, err error) {

	metric, err := findMetric(metricType, metricName, r.pgx)

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
	metricBody, err = findMetric(body.MType, body.ID, r.pgx)

	return
}

func (r *PgRepository) UpdateCounter(metricName string, delta int64) error {
	metric := models.MetricsBody{
		ID:    metricName,
		MType: models.Counter,
		Delta: delta,
	}
	_, err := findMetric(models.Counter, metricName, r.pgx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return insert(metric, r.pgx)
		}
		return err
	}

	return updateCounter(metric, r.pgx)
}

func (r *PgRepository) UpdateGauge(metricName string, value float64) error {
	metric := models.MetricsBody{
		ID:    metricName,
		MType: models.Counter,
		Value: value,
	}
	_, err := findMetric(models.Gauge, metricName, r.pgx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return insert(metric, r.pgx)
		}
		return err
	}

	return updateGauge(metric, r.pgx)
}

func (r *PgRepository) Updates(metrics []models.MetricsBody) error {
	tx, err := r.pgx.Begin(context.TODO())
	if err != nil {
		return err
	}

	for _, v := range metrics {
		_, err := findMetric(v.MType, v.ID, tx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if err := insert(v, tx); err != nil {
					tx.Rollback(context.TODO())
					return err
				}
			} else {
				tx.Rollback(context.TODO())
				return err
			}
		}

		switch v.MType {
		case models.Counter:
			if err := updateCounter(v, tx); err != nil {
				tx.Rollback(context.TODO())
				return err
			}
		case models.Gauge:
			if err := updateGauge(v, tx); err != nil {
				tx.Rollback(context.TODO())
				return err
			}
		default:
			return errors.New("wrong type")
		}
	}

	return tx.Commit(context.TODO())
}

func (r *PgRepository) Close() {
	r.pgx.Close(context.Background())
}

func (r *PgRepository) Ping(ctx context.Context) error {
	return r.pgx.Ping(ctx)
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
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		fmt.Println("migration: no change")
	}

	return nil
}

func findMetric(metricType models.MType, metricName string, conn conn) (metricBody models.Metrics, err error) {
	metricBody = models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	err = conn.QueryRow(context.TODO(), `
	SELECT delta, value
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": metricType, "metricName": metricName}).Scan(&metricBody.Delta, &metricBody.Value)

	return
}

func insert(metric models.MetricsBody, conn conn) error {
	_, err := conn.Exec(context.TODO(), `
	INSERT INTO metrics (id, m_type, delta, value)
	VALUES (@metricName, @metricType, @delta, @value)
	`, pgx.NamedArgs{
		"metricType": metric.MType,
		"metricName": metric.ID,
		"delta":      metric.Delta,
		"value":      metric.Value,
	})
	return err
}

func updateCounter(metric models.MetricsBody, conn conn) error {
	_, err := conn.Exec(context.TODO(), `
	UPDATE metrics
	SET delta = delta + @delta
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"delta": metric.Delta, "metricType": metric.MType, "metricName": metric.ID})

	return err
}

func updateGauge(metric models.MetricsBody, conn conn) error {
	_, err := conn.Exec(context.TODO(), `
	UPDATE metrics
	SET value = @value
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"value": metric.Value, "metricType": metric.MType, "metricName": metric.ID})

	return err
}
