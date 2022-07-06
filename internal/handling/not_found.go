package handling

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServeNotFoundResponse(c *gin.Context) {
	c.HTML(http.StatusNotFound, "not_found.html", struct{}{})
}
