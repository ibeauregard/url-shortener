#!/bin/sh

sqlite3 -batch "$PWD/db/url-mappings.db" <"$PWD/db/init-db.sql"
./url_shortener
