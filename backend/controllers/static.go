package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"super-supply-chain/configs"
	"super-supply-chain/utils"
)

func LoadStatic(c *gin.Engine) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !configs.IsDev() {
		root := filepath.Dir(pwd)
		staticPath := filepath.Join(root, "frontend/dist")
		utils.Logger.Info("Load static path: " + staticPath)
		c.Static("/super-supply-chain", staticPath)

		docsPath := filepath.Join(root, "website/build")
		utils.Logger.Info("Load docs path: " + docsPath)
		registerDocs(c, docsPath)
	}
}

// registerDocs serves the Docusaurus build at the host path /docs.
// The site uses trailingSlash, so each page is a directory containing index.html.
func registerDocs(engine *gin.Engine, docsPath string) {
	root, err := filepath.Abs(docsPath)
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.Dir(root))
	handler := func(c *gin.Context) {
		serveDocs(c, root, fileServer)
	}
	engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/")
	})
	engine.GET("/docs/*filepath", handler)
	engine.HEAD("/docs/*filepath", handler)
}

func serveDocs(c *gin.Context, root string, fileServer http.Handler) {
	rel := pathClean(c.Param("filepath"))
	if rel == "" {
		c.File(filepath.Join(root, "index.html"))
		return
	}
	full := filepath.Join(root, rel)
	if !isInside(root, full) {
		c.Status(http.StatusForbidden)
		return
	}
	info, err := os.Stat(full)
	if err != nil {
		if os.IsNotExist(err) {
			writeDocsNotFound(c, root)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		if !strings.HasSuffix(c.Request.URL.Path, "/") {
			c.Redirect(http.StatusMovedPermanently, c.Request.URL.Path+"/")
			return
		}
		index := filepath.Join(full, "index.html")
		if _, statErr := os.Stat(index); statErr != nil {
			writeDocsNotFound(c, root)
			return
		}
		c.File(index)
		return
	}
	http.StripPrefix("/docs", fileServer).ServeHTTP(c.Writer, c.Request)
}

func writeDocsNotFound(c *gin.Context, root string) {
	notFound := filepath.Join(root, "404.html")
	body, err := os.ReadFile(notFound)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusNotFound, "text/html; charset=utf-8", body)
}

func pathClean(rel string) string {
	rel = strings.TrimPrefix(rel, "/")
	cleaned := filepath.Clean(rel)
	if cleaned == "." {
		return ""
	}
	return cleaned
}

func isInside(root, full string) bool {
	root = filepath.Clean(root)
	full = filepath.Clean(full)
	if full == root {
		return true
	}
	sep := string(os.PathSeparator)
	return strings.HasPrefix(full, root+sep)
}
