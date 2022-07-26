package handling

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleGetFromKey(repo RepoAdapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		handle(c, repo)
	}
}

func handle(ctx *gin.Context, repo RepoAdapter) {
	key := ctx.Param("key")
	longUrl, found := repo.getLongUrl(key)
	if found {
		ctx.Redirect(http.StatusMovedPermanently, longUrl)
	} else {
		ServeNotFoundResponse(ctx)
	}
}
