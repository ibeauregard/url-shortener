package handling

import (
	"fmt"
	repo "github.com/ibeauregard/url-shortener/internal/repository"
)

type Repository interface {
	Close() error
	FindByLongUrl(longUrl string) (*repo.MappingModel, error)
	FindByKey(key string) (*repo.MappingModel, error)
	Create(mapping *repo.MappingModel) error
	GetLastId() uint
}

type repoAdapter struct {
	Repository
}

func NewRepoAdapter(r Repository) *repoAdapter {
	return &repoAdapter{r}
}

func (r *repoAdapter) Close() error {
	err := r.Repository.Close()
	if err != nil {
		return fmt.Errorf("handling.Close: %w", err)
	}
	return nil
}

func (r *repoAdapter) getKey(longUrl string) (key string, found bool) {
	mapping, err := r.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func (r *repoAdapter) getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := r.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func (r *repoAdapter) addMapping(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl, r.getNextDatabaseId())
	err = r.Create(&repo.MappingModel{
		Key:     key,
		LongUrl: longUrl,
	})
	if err != nil {
		return "", fmt.Errorf("handling.addMapping: %w", err)
	}
	return key, nil
}

func (r *repoAdapter) getNextDatabaseId() uint {
	return r.GetLastId() + 1
}
