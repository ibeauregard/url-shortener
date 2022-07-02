package main

func getKeyFromDb(longUrl string) (shortUrl string, found bool) {
	row := keyFromLongUrlStmt.QueryRow(longUrl)
	var key string
	if err := row.Scan(&key); err != nil {
		return "", false
	}
	return key, true
}

func addMappingToDb(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl)
	_, err = insertStmt.Exec(key, longUrl)
	if err != nil {
		return "", err
	}
	return key, nil
}

func getNextDatabaseId() uint {
	row := lastIdStmt.QueryRow()
	var lastId uint
	if err := row.Scan(&lastId); err != nil {
		lastId = 0
	}
	return lastId + 1
}
