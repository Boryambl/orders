package storage

import (
	"context"
	"log"
	"orders/sql"
	"sync"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

type SQLRepository struct {
	URL           string
	Pool          *pgxpool.Pool
	DBLock        *sync.Mutex
	SchemaChanged bool
}

var sqlRepo *SQLRepository

func setSQLRepository(rep *SQLRepository) {
	sqlRepo = rep
}

func SQLRepo() *SQLRepository {
	return sqlRepo
}

func (f *SQLRepository) Connect() error {
	config, err := pgxpool.ParseConfig(f.URL)
	if err != nil {
		log.Printf("Failed to parse PG connection URL %v", err)
		return err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err == nil {
		f.Pool = pool
		f.DBLock = &sync.Mutex{}
	} else {
		log.Printf("Failed to create connection pool %v", err)
	}
	return err
}

func (f *SQLRepository) Close() {
	if f == nil {
		return
	}
	if f.Pool != nil {
		f.Pool.Close()
		f.Pool = nil
	}
}

func InitSQL(path string) error {
	sqlRepo := &SQLRepository{
		URL: path,
	}
	err := sqlRepo.Connect()
	if err != nil {
		return err
	}
	setSQLRepository(sqlRepo)
	ctx := context.Background()
	tx, err := sqlRepo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("failed to start transaction %v", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, sql.Up)
	if err != nil {
		log.Printf("failed to apply initialization script %v", err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
