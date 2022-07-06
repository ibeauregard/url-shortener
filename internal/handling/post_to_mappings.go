package handling

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
)

var AppHost = os.Getenv("APP_HOST")

const AppScheme = "http"

type concretePostHandler struct {
	ctx *gin.Context
}

func HandlePostToMappings(repo RepoProxy) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		(&concretePostHandler{ctx}).handle(repo)
	}
}

func (handler *concretePostHandler) handle(repo RepoProxy) {
	normalizedUrl, err := handler.getNormalizedUrl()
	if err != nil {
		return
	}
	if normalizedUrl.Host == AppHost {
		handler.ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": fmt.Sprintf("Domain %s is blacklisted", AppHost),
		})
		return
	}
	normalizedUrlAsString := normalizedUrl.String()
	key, found := repo.getKey(normalizedUrlAsString)
	if found {
		handler.ctx.JSON(http.StatusOK, getSuccessResponseBody(normalizedUrlAsString, getShortUrl(key)))
		return
	}
	key, err = repo.addMapping(normalizedUrlAsString)
	if err != nil {
		handler.ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error: could not insert into database",
		})
		return
	}
	handler.ctx.JSON(http.StatusCreated, getSuccessResponseBody(normalizedUrlAsString, getShortUrl(key)))
}

func (handler *concretePostHandler) getNormalizedUrl() (*url.URL, error) {
	var payload struct {
		LongUrl string `json:"longUrl" binding:"required"`
	}
	err := handler.ctx.ShouldBindJSON(&payload)
	if err != nil {
		handler.ctx.JSON(http.StatusBadRequest, gin.H{"message": "Expected JSON with non-empty 'longUrl' attribute"})
		return nil, err
	}
	normalizedUrl, err := normalize(payload.LongUrl)
	if err != nil {
		handler.ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Not a valid URL"})
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
