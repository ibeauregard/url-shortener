package main

import repo "github.com/ibeauregard/url-shortener/internal/repository"

func getKey(longUrl string) (key string, found bool) {
	mapping, err := repository.FindByLongUrl(longUrl)
	if err != nil {
		return "", false
	}
	return mapping.Key, true
}

func getLongUrl(key string) (longUrl string, found bool) {
	mapping, err := repository.FindByKey(key)
	if err != nil {
		return "", false
	}
	return mapping.LongUrl, true
}

func addMapping(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl)
	err = repository.Create(&repo.MappingModel{
		Key:     key,
		LongUrl: longUrl,
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func getNextDatabaseId() uint {
	return repository.GetLastId() + 1
}
