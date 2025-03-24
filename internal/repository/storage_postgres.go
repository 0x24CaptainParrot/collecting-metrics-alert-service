package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type StoragePostgres struct {
	db *sql.DB
}

func NewStoragePostgres(db *sql.DB) *StoragePostgres {
	if db == nil {
		log.Println("Warning: database is not configured. Using in-memory storage.")
		return nil
	}

	st := &StoragePostgres{db: db}
	if err := st.CreateTables(); err != nil {
		log.Fatalf("Error creating tables in db.")
	}
	return st
}

func (st *StoragePostgres) CreateTables() error {
	return DoDBWithRetry(func() error {
		_, err := st.db.Exec(`
			CREATE TABLE IF NOT EXISTS metrics (
				id TEXT PRIMARY KEY,
				type TEXT NOT NULL CHECK (type IN ('gauge', 'counter')),
				value DOUBLE PRECISION,
				delta BIGINT
			);`)
		return err
	})
}

func (sp *StoragePostgres) UpdateMetricValue(ctx context.Context, name string, value interface{}) error {
	return DoDBWithRetry(func() error {
		switch v := value.(type) {
		case float64:
			_, err := sp.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) 
									VALUES ($1, 'gauge', $2) ON CONFLICT (id) 
								DO UPDATE SET value = $2`, name, v)
			return err
		case int64:
			_, err := sp.db.ExecContext(ctx, `INSERT INTO metrics (id, type, delta)
									VALUES ($1, 'counter', $2) ON CONFLICT (id)
								DO UPDATE SET delta = metrics.delta + $2`, name, v)
			return err
		default:
			return fmt.Errorf("invalid metric type was given: %v", v)
		}
	})
}

func (sp *StoragePostgres) SaveLoadMetrics(filePath string, operation string) error {
	return storage.NewMemStorage().SaveLoadMetrics(filePath, operation)
}

func (sp *StoragePostgres) GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error) {
	var res interface{}
	err := DoDBWithRetry(func() error {
		var mType string
		var gaugeValue sql.NullFloat64
		var counterValue sql.NullInt64

		row := sp.db.QueryRowContext(ctx, "SELECT type, value, delta FROM metrics WHERE id = $1", name)
		scanErr := row.Scan(&mType, &gaugeValue, &counterValue)
		if scanErr != nil {
			return scanErr
		}

		if mType == "gauge" && gaugeValue.Valid {
			res = gaugeValue.Float64
			return nil
		} else if mType == "counter" && counterValue.Valid {
			res = counterValue.Int64
			return nil
		}

		return fmt.Errorf("metric not found")
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sp *StoragePostgres) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	err := DoDBWithRetry(func() error {
		rows, err := sp.db.QueryContext(ctx, "SELECT id, type, value, delta FROM metrics")
		if err != nil {
			log.Printf("Error querying metrics: %v", err)
			return err
		}
		defer rows.Close()

		locMetrics := make(map[string]interface{})
		for rows.Next() {
			var id, mType string
			var gaugeValue sql.NullFloat64
			var counterValue sql.NullInt64

			err := rows.Scan(&id, &mType, &gaugeValue, &counterValue)
			if err != nil {
				log.Printf("Error scanning metric row: %v", err)
				return err
			}

			if mType == "gauge" && gaugeValue.Valid {
				metrics[id] = gaugeValue.Float64
			} else if mType == "counter" && counterValue.Valid {
				metrics[id] = counterValue.Int64
			}
		}

		if rowsErr := rows.Err(); rowsErr != nil {
			log.Printf("Error iterating over metric rows: %v", err)
			return rowsErr
		}

		metrics = locMetrics
		return nil
	})
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (sp *StoragePostgres) Ping() error {
	return DoDBWithRetry(func() error {
		return sp.db.Ping()
	})
}

func (sp *StoragePostgres) DB() *sql.DB {
	return sp.db
}

func DoDBWithRetry(fn func() error) error {
	var backoffs = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var lastErr error
	for i := 0; i < len(backoffs)+1; i++ {
		err := fn()
		if err != nil {
			if IsRetriableDBErr(err) && i < len(backoffs) {
				lastErr = err
				time.Sleep(backoffs[i])
				continue
			}
			return err
		}
		return nil
	}
	return lastErr
}

func IsRetriableDBErr(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	return false
}
