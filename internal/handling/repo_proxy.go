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

type repoProxy struct {
	r Repository
}

func NewRepoProxy(r Repository) RepoProxy {
	return &repoProxy{r: r}
}

func (proxy *repoProxy) Close() error {
	err := proxy.r.Close()
	if err != nil {
		return fmt.Errorf("handling.Close: %w", err)
	}
	return nil
}

func (proxy *repoProxy) getKey(longUrl string) (key string, found bool) {
	mapping, err := proxy.r.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func (proxy *repoProxy) getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := proxy.r.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func (proxy *repoProxy) addMapping(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl, proxy.getNextDatabaseId())
	err = proxy.r.Create(&repo.MappingModel{
		Key:     key,
		LongUrl: longUrl,
	})
	if err != nil {
		return "", fmt.Errorf("handling.addMapping: %w", err)
	}
	return key, nil
}

func (proxy *repoProxy) getNextDatabaseId() uint {
	return proxy.r.GetLastId() + 1
}
