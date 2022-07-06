package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ibeauregard/url-shortener/internal/handling"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	var r = gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Print(err)
	}
	repo, _ := handling.NewRepoProxy("db/data/url-mappings.db")
	defer repo.Close()
	performRouting(r, repo)
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
