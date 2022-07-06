package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func serveNotFoundResponse(c *gin.Context) {
	c.HTML(http.StatusNotFound, "not_found.html", struct{}{})
}
