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
		log.Print(err)
	}
	repo, err := sqlite.NewRepository("db/data/url-mappings.db", sql.Open)
	if err != nil {
		log.Panic(err)
	}
	repoProxy := handling.NewRepoProxy(repo)
	defer repoProxy.Close()
	performRouting(r, repoProxy)
	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func performRouting(r *gin.Engine, repo handling.RepoProxy) {
	r.LoadHTMLFiles("templates/not_found.html")
	r.Static("/static", "./static")
	r.POST("/api/mappings", handling.HandlePostToMappings(repo))
	r.GET("/:key", handling.HandleGetFromKey(repo))
	r.GET("/", handling.ServeNotFoundResponse)
}
