package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	var r = gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Print(err)
	}
	repo, _ := newRepoProxy("db/data/url-mappings.db")
	defer repo.close()
	performRouting(r, repo)
	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func performRouting(r *gin.Engine, repo repoProxy) {
	r.LoadHTMLFiles("templates/not_found.html")
	r.Static("/static", "./static")
	r.POST("/api/mappings", handlePostToMappings(repo))
	r.GET("/:key", handleGetFromKey(repo))
	r.GET("/", serveNotFoundResponse)
}
