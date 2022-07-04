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
	r.LoadHTMLFiles("templates/not_found.html")
	r.Static("/static", "./static")
	r.POST("/api/mappings", handlePostToMappings)
	r.GET("/:key", handleGetFromKey)
	r.GET("/", serveNotFoundResponse)
	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}
