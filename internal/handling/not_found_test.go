package handling

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeNotFoundResponse(t *testing.T) {
	testName := "serveNotFoundResponse(c *gin.Context)"
	t.Run(testName, func(t *testing.T) {
		r := gin.Default()
		r.GET("/not-found", ServeNotFoundResponse)
		req, _ := http.NewRequest("GET", "/not-found", nil)
		r.LoadHTMLFiles("../templates/not_found.html")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.EqualValues(t, http.StatusNotFound, w.Code)
	})
}
