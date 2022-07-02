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
	r.POST("/api/mappings", handlePostToMappings)
	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}
