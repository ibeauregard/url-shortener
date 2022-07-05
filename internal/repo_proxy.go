package main

import (
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"github.com/ibeauregard/url-shortener/internal/repository/sqlite"
	"log"
)

type repoProxy interface {
	close()
	getKey(longUrl string) (key string, found bool)
	getLongUrl(key string) (longUrl string, found bool)
	addMapping(longUrl string) (shortUrl string, err error)
}

type repoProxyImpl struct {
	r repo.Repository
}

func newRepoProxy(dataSourceName string) (repoProxy, error) {
	repository, err := sqlite.NewRepository(dataSourceName)
	if err != nil {
		return nil, err
	}
	return &repoProxyImpl{r: repository}, nil
}

func (proxy *repoProxyImpl) close() {
	err := proxy.r.Close
	if err != nil {
		log.Printf("Unable to close repo %v", proxy.r)
	}
}

func (proxy *repoProxyImpl) getKey(longUrl string) (key string, found bool) {
	mapping, err := proxy.r.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func (proxy *repoProxyImpl) getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := proxy.r.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func (proxy *repoProxyImpl) addMapping(longUrl string) (shortUrl string, err error) {
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

func (proxy *repoProxyImpl) getNextDatabaseId() uint {
	return proxy.r.GetLastId() + 1
}
