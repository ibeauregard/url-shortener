package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/ibeauregard/url-shortener/internal/handling"
	"github.com/ibeauregard/url-shortener/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	var r = gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Printf("main: %v", err)
	}
	repo, err := sqlite.NewRepository("db/data/url-mappings.db", sql.Open)
	if err != nil {
		log.Panicf("main: %v", err)
	}
	repoProxy := handling.NewRepoProxy(repo)
	defer func() {
		if err := repoProxy.Close(); err != nil {
			log.Printf("main: %v", err)
		}
	}()
	handling.PerformRouting(r, repoProxy)
	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}
