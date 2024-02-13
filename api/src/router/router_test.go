package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"api/router"
)

type mockUserAndEntryController struct{}

func (m *mockUserAndEntryController) HandleUserAndEntry(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TestNewRouter(t *testing.T) {
	ctrl := &mockUserAndEntryController{}
	router := router.NewRouter(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!", rec.Body.String())
}
