package sqlite

import (
	"database/sql"
	"fmt"
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"golang.org/x/net/context"
	"log"
	"time"
)

type sqlDb interface {
	Close() error
	Ping() error
	QueryRowContext(context.Context, string, ...any) *sql.Row
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

type repository struct {
	db sqlDb
}

type sqlOpener func(string, string) (*sql.DB, error)

func NewRepository(dataSourceName string, sqlOpen sqlOpener) (*repository, error) {
	db, err := sqlOpen("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewRepository: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewRepository: %w", err)
	}
	return &repository{db}, nil
}

func (r *repository) Close() error {
	err := r.db.Close()
	if err != nil {
		return fmt.Errorf("sqlite.Close: %w", err)
	}
	return nil
}

func (r *repository) FindByLongUrl(longUrl string) (*repo.MappingModel, error) {
	mapping := new(repo.MappingModel)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, "SELECT id, key, long_url FROM mappings WHERE long_url=?", longUrl).
		Scan(&mapping.ID, &mapping.Key, &mapping.LongUrl)
	if err != nil {
		return nil, fmt.Errorf("sqlite.FindByLongUrl: %w", err)
	}
	return mapping, nil
}

func (r *repository) FindByKey(key string) (*repo.MappingModel, error) {
	mapping := new(repo.MappingModel)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, "SELECT id, key, long_url FROM mappings WHERE key = ?", key).
		Scan(&mapping.ID, &mapping.Key, &mapping.LongUrl)
	if err != nil {
		return nil, fmt.Errorf("sqlite.FindByKey: %w", err)
	}
	return mapping, nil
}

func (r *repository) Create(mapping *repo.MappingModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO mappings (key, long_url) VALUES(?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("sqlite.Create: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("sqlite.Create: %v", err)
		}
	}()
	_, err = stmt.ExecContext(ctx, mapping.Key, mapping.LongUrl)
	return err
}

func (r *repository) GetLastId() uint {
	var lastId uint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, "SELECT seq FROM sqlite_sequence WHERE name='mappings'").Scan(&lastId)
	if err != nil {
		lastId = 0
	}
	return lastId
}
