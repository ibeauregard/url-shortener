package handling

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetFromKey(t *testing.T) {
	assert.NotNil(t, HandleGetFromKey(&concreteRepoProxy{}))
}

type repoProxyMock struct {
	RepoProxy
	outputStr         string
	outputFoundStatus bool
	outputError       error
}

func (m *repoProxyMock) getLongUrl(_ string) (string, bool) {
	return m.outputStr, m.outputFoundStatus
}

func TestHandleFound(t *testing.T) {
	dummyLongUrl := "http://foobar.com"
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Method: "GET",
	}
	ctx.AddParam("key", "my_key")
	mock := &repoProxyMock{
		outputStr:         dummyLongUrl,
		outputFoundStatus: true,
	}
	handle(ctx, mock)
	assert.EqualValues(t, http.StatusMovedPermanently, w.Code)
	assert.EqualValues(t, []string{dummyLongUrl}, w.Result().Header["Location"])
}

func TestHandleNotFound(t *testing.T) {
	mock := &repoProxyMock{
		outputStr:         "",
		outputFoundStatus: false,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/my_key_not_found", nil)
	r := gin.Default()
	r.GET("/:key", HandleGetFromKey(mock))
	r.LoadHTMLFiles("../templates/not_found.html")
	r.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusNotFound, w.Code)
}
