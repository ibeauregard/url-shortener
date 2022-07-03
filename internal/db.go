package main

import "database/sql"

func getKeyFromDb(longUrl string) (key string, found bool) {
	return getFromDb(keyFromLongUrlStmt, longUrl)
}

func getLongUrlFromDb(key string) (longUrl string, found bool) {
	return getFromDb(longUrlFromKeyStmt, key)
}

func addMappingToDb(longUrl string) (shortUrl string, err error) {
	key := generateKey(longUrl)
	_, err = insertStmt.Exec(key, longUrl)
	if err != nil {
		return "", err
	}
	return key, nil
}

func getFromDb(statement *sql.Stmt, indexKey string) (value string, found bool) {
	row := statement.QueryRow(indexKey)
	if err := row.Scan(&value); err != nil {
		return "", false
	}
	return value, true
}

func getNextDatabaseId() uint {
	row := lastIdStmt.QueryRow()
	var lastId uint
	if err := row.Scan(&lastId); err != nil {
		lastId = 0
	}
	return lastId + 1
}
