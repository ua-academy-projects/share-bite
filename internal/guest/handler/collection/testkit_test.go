package collection

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
	guest_middleware "github.com/ua-academy-projects/share-bite/internal/guest/middleware"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

var (
	internalErrMsg = "internal server error"
	validationMsg  = "request validation failed"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	binding.Validator = validator.New("binding")

	r := gin.New()
	r.Use(guest_middleware.ErrorMiddleware())
	return r
}

func performJSONRequest(t *testing.T, r http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Reader
	if body == nil {
		reqBody = bytes.NewReader(nil)
	} else {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRawJSONRequest(t *testing.T, r http.Handler, method string, path string, raw string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewBufferString(raw))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequest(r http.Handler, method string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertJSONBody(t *testing.T, want any, got string) {
	t.Helper()
	b, err := json.Marshal(want)
	require.NoError(t, err)
	require.JSONEq(t, string(b), got)
}

func withCustomerID(v any, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v != nil {
			c.Set(middleware.CtxCustomerID, v)
		}

		next(c)
	}
}

func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }
