package main

import (
	"bytes"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func serveFile(c *echo.Context) {
	path := c.Request.URL.Path
	if len(path) > 1 {
		path = path[1:len(path)]
	}

	data, err := Asset(path)
	if err != nil {
		println("Looking for " + path + " Error was " + err.Error())
		http.NotFound(c.Response, c.Request)
		return
	}
	reader := bytes.NewReader(data)
	http.ServeContent(c.Response, c.Request, filepath.Base(path), time.Now(), reader)
}

func main() {
	e := echo.New()

	e.Get("/", func(c *echo.Context) {
		c.Request.URL.Path = "/index.html"
		serveFile(c)
	})
	e.Get("/*", serveFile)
	e.Run(":" + os.Getenv("PORT"))
}
