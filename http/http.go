package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/qwqcode/qwquiver/bindata"
	"github.com/qwqcode/qwquiver/config"
	"github.com/sirupsen/logrus"
)

// Run 运行 http server
func Run() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // For dev
	}))
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	e.Use(echolog.Middleware(echolog.DefaultConfig))

	fileServer := http.FileServer(bindata.AssetFile())
	e.GET("/*", echo.WrapHandler(fileServer))

	api := e.Group("/api")
	api.GET("/query", queryHandler)
	api.GET("/conf", confHandler)
	api.GET("/chart", chartHandler)
	api.GET("/school/all", schoolAllHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Instance.Port)))
}
