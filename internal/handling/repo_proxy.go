package handling

import (
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"log"
)

type RepoProxy interface {
	Close() error
	getKey(longUrl string) (key string, found bool)
	getLongUrl(key string) (longUrl string, found bool)
	addMapping(longUrl string) (shortUrl string, err error)
}

type concreteRepoProxy struct {
	r repo.Repository
}

func NewRepoProxy(r repo.Repository) RepoProxy {
	return &concreteRepoProxy{r: r}
}

func (proxy *concreteRepoProxy) Close() error {
	err := proxy.r.Close()
	if err != nil {
		log.Printf("Unable to close repo %v", proxy.r)
	}
	return err
}

func (proxy *concreteRepoProxy) getKey(longUrl string) (key string, found bool) {
	mapping, err := proxy.r.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func (proxy *concreteRepoProxy) getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := proxy.r.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func (proxy *concreteRepoProxy) addMapping(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl, proxy.getNextDatabaseId())
	err = proxy.r.Create(&repo.MappingModel{
		Key:     key,
		LongUrl: longUrl,
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (proxy *concreteRepoProxy) getNextDatabaseId() uint {
	return proxy.r.GetLastId() + 1
}
