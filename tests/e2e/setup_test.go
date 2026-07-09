package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/handlers"
	"github.com/chop1k/medods-test/internal/routes"
)

// newTestServer boots the real Gin router (same wiring as main.go) behind an
// httptest.Server so tests exercise the full HTTP stack - routing, binding
// and validation - over a real HTTP client.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/v1")
	routes.RegisterTemplateRoutes(v1, handlers.NewTemplateHandler())
	routes.RegisterTaskRoutes(v1, handlers.NewTaskHandler())
	routes.RegisterTagRoutes(v1, handlers.NewTagHandler())

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	return srv
}

// doJSON performs an HTTP request with an optional JSON body and decodes the
// JSON response (if any) into out.
func doJSON(t *testing.T, client *http.Client, method, url string, body any, out any) *http.Response {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	if out != nil && resp.ContentLength != 0 {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
	}

	return resp
}
