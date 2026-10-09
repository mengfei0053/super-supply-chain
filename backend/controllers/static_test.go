package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterDocs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "index.html"), "home")
	mustWrite(t, filepath.Join(dir, "intro", "index.html"), "intro")
	mustWrite(t, filepath.Join(dir, "assets", "app.js"), "js")
	mustWrite(t, filepath.Join(dir, "404.html"), "missing")

	r := gin.New()
	registerDocs(r, dir)

	t.Run("root redirect", func(t *testing.T) {
		rec := get(r, "/docs")
		if rec.Code != http.StatusMovedPermanently {
			t.Fatalf("status %d", rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != "/docs/" {
			t.Fatalf("location %q", loc)
		}
	})

	t.Run("home", func(t *testing.T) {
		rec := get(r, "/docs/")
		if rec.Code != http.StatusOK || rec.Body.String() != "home" {
			t.Fatalf("got %d %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("page redirect and index", func(t *testing.T) {
		rec := get(r, "/docs/intro")
		if rec.Code != http.StatusMovedPermanently {
			t.Fatalf("status %d", rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != "/docs/intro/" {
			t.Fatalf("location %q", loc)
		}
		rec = get(r, "/docs/intro/")
		if rec.Code != http.StatusOK || rec.Body.String() != "intro" {
			t.Fatalf("got %d %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("asset", func(t *testing.T) {
		rec := get(r, "/docs/assets/app.js")
		if rec.Code != http.StatusOK || rec.Body.String() != "js" {
			t.Fatalf("got %d %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("missing", func(t *testing.T) {
		rec := get(r, "/docs/nope/")
		if rec.Code != http.StatusNotFound || rec.Body.String() != "missing" {
			t.Fatalf("got %d %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("traversal", func(t *testing.T) {
		outside := filepath.Join(filepath.Dir(dir), "secret.txt")
		mustWrite(t, outside, "secret")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/docs/x", nil)
		c.Params = gin.Params{{Key: "filepath", Value: "/../secret.txt"}}
		serveDocs(c, dir, http.FileServer(http.Dir(dir)))
		if c.Writer.Status() != http.StatusForbidden || strings.Contains(w.Body.String(), "secret") {
			t.Fatalf("got %d %q", c.Writer.Status(), w.Body.String())
		}
	})
}

func get(r http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
