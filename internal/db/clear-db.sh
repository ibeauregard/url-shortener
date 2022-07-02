#!/bin/sh

sqlite3 -batch "$PWD/db/data/url-mappings.db" <"$PWD/db/clear-db.sql"
