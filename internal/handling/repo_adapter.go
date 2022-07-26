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
	r Repository
}

func NewRepoAdapter(r Repository) *repoAdapter {
	return &repoAdapter{r: r}
}

func (adapter *repoAdapter) Close() error {
	err := adapter.r.Close()
	if err != nil {
		return fmt.Errorf("handling.Close: %w", err)
	}
	return nil
}

func (adapter *repoAdapter) getKey(longUrl string) (key string, found bool) {
	mapping, err := adapter.r.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func (adapter *repoAdapter) getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := adapter.r.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func (adapter *repoAdapter) addMapping(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl, adapter.getNextDatabaseId())
	err = adapter.r.Create(&repo.MappingModel{
		Key:     key,
		LongUrl: longUrl,
	})
	if err != nil {
		return "", fmt.Errorf("handling.addMapping: %w", err)
	}
	return key, nil
}

func (adapter *repoAdapter) getNextDatabaseId() uint {
	return adapter.r.GetLastId() + 1
}
