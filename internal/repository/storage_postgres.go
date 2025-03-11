package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
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
	_, err := st.db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL CHECK (type IN ('gauge', 'counter')),
			value DOUBLE PRECISION,
			delta BIGINT
		);`)
	return err
}

func (sp *StoragePostgres) UpdateGauge(name string, value float64) error {
	_, err := sp.db.Exec("INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2) "+
		"ON CONFLICT (id) DO UPDATE SET value = $2", name, value)

	return err
}

func (sp *StoragePostgres) UpdateCounter(name string, value int64) error {
	_, err := sp.db.Exec("INSERT INTO metrics (id, type, delta) VALUES ($1, 'counter', $2) "+
		"ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + $2", name, value)

	return err
}

func (sp *StoragePostgres) GetMetric(name string, metricType storage.MetricType) (interface{}, error) {
	var mType string
	var gaugeValue sql.NullFloat64
	var counterValue sql.NullInt64

	row := sp.db.QueryRow("SELECT type, value, delta FROM metrics WHERE id = $1", name)
	err := row.Scan(&mType, &gaugeValue, &counterValue)
	if err != nil {
		return nil, err
	}

	if mType == "gauge" && gaugeValue.Valid {
		return gaugeValue.Float64, nil
	} else if mType == "counter" && counterValue.Valid {
		return counterValue.Int64, nil
	}

	return nil, fmt.Errorf("metric not found")
}

func (sp *StoragePostgres) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	rows, err := sp.db.Query("SELECT id, type, value, delta FROM metrics")
	if err != nil {
		log.Printf("Error querying metrics: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id, mType string
		var gaugeValue sql.NullFloat64
		var counterValue sql.NullInt64

		err := rows.Scan(&id, &mType, &gaugeValue, &counterValue)
		if err != nil {
			log.Printf("Error scanning metric row: %v", err)
			continue
		}

		if mType == "gauge" && gaugeValue.Valid {
			metrics[id] = gaugeValue.Float64
		} else if mType == "counter" && counterValue.Valid {
			metrics[id] = counterValue.Int64
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over metric rows: %v", err)
	}

	return metrics
}

func (sp *StoragePostgres) Ping() error {
	return sp.db.Ping()
}

func (sp *StoragePostgres) DB() *sql.DB {
	return sp.db
}

func (sp *StoragePostgres) SaveMetricsToFile(filePath string) error {
	return storage.NewMemStorage().SaveMetricsToFile(filePath)
}

func (sp *StoragePostgres) LoadMetricsFromFile(filePath string) error {
	return storage.NewMemStorage().LoadMetricsFromFile(filePath)
}
