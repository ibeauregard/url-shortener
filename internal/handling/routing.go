package handling

import "github.com/gin-gonic/gin"

type RepoProxy interface {
	Close() error
	getKey(longUrl string) (key string, found bool)
	getLongUrl(key string) (longUrl string, found bool)
	addMapping(longUrl string) (shortUrl string, err error)
}

func PerformRouting(r *gin.Engine, repo RepoProxy) {
	r.LoadHTMLFiles("templates/not_found.html")
	r.Static("/static", "./static")
	r.POST("/api/mappings", HandlePostToMappings(repo))
	r.GET("/:key", HandleGetFromKey(repo))
	r.GET("/", ServeNotFoundResponse)
}
