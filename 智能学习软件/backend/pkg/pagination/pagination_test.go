package pagination_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"smart-learning/pkg/pagination"
)

func setupCtx(query string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/?"+query, nil)
	c.Request = req
	return c
}

func TestParse_Defaults(t *testing.T) {
	c := setupCtx("")
	p := pagination.Parse(c)
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.PageSize)
	assert.Equal(t, 0, p.Offset())
	assert.Equal(t, 20, p.Limit())
}

func TestParse_CustomValues(t *testing.T) {
	c := setupCtx("page=3&page_size=50")
	p := pagination.Parse(c)
	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 50, p.PageSize)
	assert.Equal(t, 100, p.Offset())
	assert.Equal(t, 50, p.Limit())
}

func TestParse_Clamping(t *testing.T) {
	c := setupCtx("page=0&page_size=-1")
	p := pagination.Parse(c)
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.PageSize)
}

func TestParse_MaxPageSize(t *testing.T) {
	c := setupCtx("page_size=500")
	p := pagination.Parse(c)
	assert.Equal(t, 100, p.PageSize)
}

func TestParse_Invalid(t *testing.T) {
	c := setupCtx("page=abc&page_size=xyz")
	p := pagination.Parse(c)
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.PageSize)
}