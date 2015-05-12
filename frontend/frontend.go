package main

import (
	"bytes"
	"github.com/abates/vpanel"
	"github.com/labstack/echo"
	"net/http"
	"path/filepath"
	"time"
)

func serveFile(c *echo.Context) {
	path := c.Request.URL.Path
	if len(path) > 1 {
		path = path[1:len(path)]
	}

	data, err := vpanel.Asset(path)
	if err != nil {
		http.NotFound(c.Response, c.Request)
		return
	}
	reader := bytes.NewReader(data)
	http.ServeContent(c.Response, c.Request, filepath.Base(path), time.Now(), reader)
}

type errorInfo struct {
	Message string `json:"message"`
}

func renderJSON(c *echo.Context, data interface{}, err error) error {
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func bindAndValidate(c *echo.Context, receiver vpanel.Validatable) error {
	if err := c.Bind(receiver); err != nil {
		return err
	}
	return receiver.Validate()
}

func main() {
	e := echo.New()

	e.Get("/", func(c *echo.Context) {
		c.Request.URL.Path = "/index.html"
		serveFile(c)
	})
	e.Get("/*", serveFile)

	a := e.Group("/api")

	a.Get("/host", func(c *echo.Context) error {
		data, err := vpanel.Host.Stats()
		return renderJSON(c, data, err)
	})

	a.Get("/containers/templates", func(c *echo.Context) error {
		data, err := vpanel.Containers.Templates()
		return renderJSON(c, data, err)
	})

	a.Get("/containers", func(c *echo.Context) error {
		data, err := vpanel.Containers.All()
		return renderJSON(c, data, err)
	})

	a.Post("/containers", func(c *echo.Context) error {
		metadata := vpanel.NewContainerMetadata()
		if err := bindAndValidate(c, metadata); err != nil {
			return err
		}
		container, err := vpanel.Containers.Create(metadata)
		return renderJSON(c, container, err)
	})

	a.HTTPErrorHandler(func(code int, err error, c *echo.Context) {
		vpanel.Logger.Warnf("Failed to process %s %s: %s", c.Request.Method, c.Request.URL.Path, err.Error())
		http.Error(c.Response, err.Error(), code)
	})

	e.Run(":" + vpanel.Config["listenPort"])
}
