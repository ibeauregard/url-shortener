#!/bin/sh

sqlite3 -batch "$PWD/db/data/url-mappings.db" <"$PWD/db/init-db.sql"
mv "$PWD/coverage.txt" "$PWD/tests/unit/"
./url_shortener
