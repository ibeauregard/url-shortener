package sqlite

import (
	"database/sql"
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"golang.org/x/net/context"
	"log"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dialect, dataSourceName string) (repo.Repository, error) {
	db, err := sql.Open(dialect, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) Close() {
	err := r.db.Close()
	if err != nil {
		log.Printf("Unable to close database %v", r.db)
	}
}

func (r *repository) FindByLongUrl(longUrl string) (*repo.MappingModel, error) {
	mapping := new(repo.MappingModel)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, "SELECT id, key, long_url FROM mappings WHERE long_url=?", longUrl).
		Scan(&mapping.ID, &mapping.Key, &mapping.LongUrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return mapping, nil
}

func (r *repository) Create(mapping *repo.MappingModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO mappings (key, long_url) VALUES(?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Unable to close statement %v", stmt)
		}
	}()
	_, err = stmt.ExecContext(ctx, mapping.Key, mapping.LongUrl)
	return err
}

func (r *repository) GetLastId() uint {
	var lastId uint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, "SELECT seq FROM sqlite_sequence WHERE name='mappings'").
		Scan(&lastId)
	if err != nil {
		lastId = 0
	}
	return lastId
}
