package handling

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerformRouting(t *testing.T) {
	r := gin.Default()
	PerformRouting(r, &repoProxy{})
	assert.NotNil(t, r.Routes())
}
