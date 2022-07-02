package main

import (
	"database/sql"
	"os"
)

var AppHost = os.Getenv("APP_HOST")

const AppScheme = "http"

var db, _ = sql.Open("sqlite3", "db/data/url-mappings.db")
var keyFromLongUrlStmt, _ = db.Prepare("SELECT key FROM mappings WHERE long_url=?;")
var insertStmt, _ = db.Prepare("INSERT INTO mappings (key, long_url) VALUES(?, ?)")
var lastIdStmt, _ = db.Prepare("SELECT seq FROM sqlite_sequence WHERE name='mappings'")

// This alphabet will be used to generate the paths of the shortened URLs.
// It consists of the decimal digits and of the uppercase and lowercase letters, plus some special characters.
// Characters that could cause ambiguity or generate offensive words were removed.
const alphabet = "23456789BCDFGHJKLMNPQRSTVWXYZbcdfghjkmnpqrstvwxyz-_~!$&=@"
const alphabetLength = uint(len(alphabet))
