package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func handlePostToMappings(c *gin.Context) {
	normalizedUrl, err := getNormalizedUrl(c)
	if err != nil {
		return
	}
	if normalizedUrl.Host == AppHost {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errorMessage": fmt.Sprintf("Domain %s is blacklisted", AppHost),
		})
		return
	}
	normalizedUrlAsString := normalizedUrl.String()
	key, found := getKeyFromDb(normalizedUrlAsString)
	if found {
		c.JSON(http.StatusOK, getSuccessResponseBody(normalizedUrlAsString, getShortUrl(key)))
		return
	}
	key, err = addMappingToDb(normalizedUrlAsString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMessage": "Internal Server Error: could not insert into database",
		})
		return
	}
	c.JSON(http.StatusCreated, getSuccessResponseBody(normalizedUrlAsString, getShortUrl(key)))
}

func handleGetFromKey(c *gin.Context) {
	key := c.Param("key")
	longUrl, found := getLongUrlFromDb(key)
	if found {
		c.Redirect(http.StatusMovedPermanently, longUrl)
	} else {
		c.HTML(http.StatusNotFound, "not_found.html", struct{}{})
	}
}

type RequestBody struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

func getNormalizedUrl(c *gin.Context) (*url.URL, error) {
	var payload RequestBody
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Expected JSON with non-empty 'longUrl' attribute"})
		return nil, err
	}
	normalizedUrl, err := normalize(payload.LongUrl)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errorMessage": "Not a valid URL"})
		return nil, err
	}
	return normalizedUrl, nil
}

func getShortUrl(key string) string {
	return fmt.Sprintf("%s://%s/%s", AppScheme, AppHost, key)
}

func getSuccessResponseBody(longUrl string, shortUrl string) gin.H {
	return gin.H{
		"longUrl":  longUrl,
		"shortUrl": shortUrl,
	}
}
