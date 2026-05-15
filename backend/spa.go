package main

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:public
var spaFiles embed.FS

func spaMiddleware() gin.HandlerFunc {
	sub, err := fs.Sub(spaFiles, "public")
	if err != nil {
		panic("embedded spa files not found — run 'make build' first")
	}

	fileServer := http.FileServer(http.FS(sub))

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") ||
			strings.HasPrefix(c.Request.URL.Path, "/uploads") {
			c.Next()
			return
		}

		sanitized := filepath.Clean(c.Request.URL.Path)
		if !strings.HasPrefix(sanitized, "/") {
			sanitized = "/" + sanitized
		}

		f, err := sub.Open(strings.TrimPrefix(sanitized, "/"))
		if err != nil {
			c.Request.URL.Path = "/index.html"
		} else {
			f.Close()
		}

		ext := filepath.Ext(c.Request.URL.Path)
		if ext == ".js" {
			c.Header("Content-Type", "application/javascript")
		} else if ext == ".css" {
			c.Header("Content-Type", "text/css")
		} else if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			c.Header("Content-Type", mimeType)
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
