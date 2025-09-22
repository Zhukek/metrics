package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/repository/postgresql/pgerr"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgRepository struct {
	pgx *pgx.Conn
}

type DBConnection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *PgRepository) GetList() ([]string, error) {
	var keys []string
	classifier := pgerr.NewPostgresErrorClassifier()

	rows, err := r.pgx.Query(context.TODO(), `
	SELECT id FROM metrics`)
	if err != nil {
		classification := classifier.Classify(err)
		if classification == pgerr.Retriable {
			for i := 0; i < 2; i++ {
				await := (i * 2) + 1
				time.Sleep(time.Duration(await) * time.Second)
				rows, err = r.pgx.Query(context.TODO(), `
			SELECT id FROM metrics`)
				if err == nil {
					break
				}
			}
		}
		if err != nil {
			return nil, err
		}
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

	metric, err := retryWithResult(findMetric, models.Metrics{
		MType: metricType,
		ID:    metricName,
	}, r.pgx)

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

func (r *PgRepository) GetMetricByRequest(body models.Metrics) (models.Metrics, error) {
	metricBody, err := retryWithResult(findMetric, body, r.pgx)

	if err != nil {
		return *metricBody, err
	}

	return *metricBody, nil
}

func (r *PgRepository) UpdateCounter(metricName string, delta int64) error {
	metric := models.Metrics{
		ID:    metricName,
		MType: models.Counter,
		Delta: &delta,
	}
	_, err := retryWithResult(findMetric, metric, r.pgx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return retry(insert, metric, r.pgx)
		}
		return err
	}

	return retry(updateCounter, metric, r.pgx)
}

func (r *PgRepository) UpdateGauge(metricName string, value float64) error {
	metric := models.Metrics{
		ID:    metricName,
		MType: models.Gauge,
		Value: &value,
	}
	_, err := retryWithResult(findMetric, metric, r.pgx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return retry(insert, metric, r.pgx)
		}
		return err
	}

	return retry(updateGauge, metric, r.pgx)
}

func (r *PgRepository) Updates(metrics []models.Metrics) error {
	tx, err := r.pgx.Begin(context.TODO())
	if err != nil {
		return err
	}

	for _, v := range metrics {
		if v.MType == models.Counter && v.Delta == nil {
			tx.Rollback(context.TODO())
			return errors.New("counter delta is nil")
		}
		if v.MType == models.Gauge && v.Value == nil {
			tx.Rollback(context.TODO())
			return errors.New("gauge value is nil")
		}

		_, err := retryWithResult(findMetric, v, tx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if err := retry(insert, v, tx); err != nil {
					tx.Rollback(context.TODO())
					return err
				}
				continue
			} else {
				tx.Rollback(context.TODO())
				return err
			}
		}

		switch v.MType {
		case models.Counter:
			if err := retry(updateCounter, v, tx); err != nil {
				tx.Rollback(context.TODO())
				return err
			}
		case models.Gauge:
			if err := retry(updateGauge, v, tx); err != nil {
				tx.Rollback(context.TODO())
				return err
			}
		default:
			tx.Rollback(context.TODO())
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

func findMetric(metric models.Metrics, conn DBConnection) (*models.Metrics, error) {

	err := conn.QueryRow(context.TODO(), `
	SELECT delta, value
	FROM metrics
	WHERE m_type=@metricType AND id=@metricName`,
		pgx.NamedArgs{"metricType": metric.MType, "metricName": metric.ID}).Scan(&metric.Delta, &metric.Value)

	return &metric, err
}

func insert(metric models.Metrics, conn DBConnection) error {
	query := `INSERT INTO metrics (id, m_type, `
	args := pgx.NamedArgs{
		"metricType": metric.MType,
		"metricName": metric.ID,
	}

	switch metric.MType {
	case models.Counter:
		query += `delta) VALUES (@metricName, @metricType, @delta)`
		args["delta"] = *metric.Delta
	case models.Gauge:
		query += `value) VALUES (@metricName, @metricType, @value)`
		args["value"] = *metric.Value
	default:
		return errors.New("wrong type")

	}

	_, err := conn.Exec(context.TODO(), query, args)

	return err
}

func updateCounter(metric models.Metrics, conn DBConnection) error {
	_, err := conn.Exec(context.TODO(), `
	UPDATE metrics
	SET delta = delta + @delta
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"delta": *metric.Delta, "metricType": metric.MType, "metricName": metric.ID})

	return err
}

func updateGauge(metric models.Metrics, conn DBConnection) error {

	_, err := conn.Exec(context.TODO(), `
	UPDATE metrics
	SET value = @value
	WHERE m_type = @metricType AND id = @metricName
	`, pgx.NamedArgs{"value": *metric.Value, "metricType": metric.MType, "metricName": metric.ID})

	return err
}

func retryWithResult[T any](f func(metric models.Metrics, conn DBConnection) (T, error), metric models.Metrics, conn DBConnection) (T, error) {
	intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	res, err := f(metric, conn)
	if err != nil {
		classifier := pgerr.NewPostgresErrorClassifier()
		for _, i := range intervals {
			classification := classifier.Classify(err)
			if classification != pgerr.Retriable {
				break
			}

			time.Sleep(i)
			res, err = f(metric, conn)
			if err == nil {
				return res, nil
			}
		}
		return res, err
	}
	return res, nil
}

func retry(f func(metric models.Metrics, conn DBConnection) error, metric models.Metrics, conn DBConnection) error {
	intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	err := f(metric, conn)
	if err != nil {
		classifier := pgerr.NewPostgresErrorClassifier()
		for _, i := range intervals {
			classification := classifier.Classify(err)
			if classification != pgerr.Retriable {
				break
			}

			time.Sleep(i)
			err = f(metric, conn)
			if err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}
