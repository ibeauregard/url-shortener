#!/bin/sh

sqlite3 -batch "$PWD/db/data/url-mappings.db" <"$PWD/db/init-db.sql"
./url_shortener
