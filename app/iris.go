package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/qwqcode/qwquiver/bindata"
	"github.com/qwqcode/qwquiver/config"
	"github.com/qwqcode/qwquiver/controllers/api"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/sirupsen/logrus"
)

// RunIris 运行 Iris
func RunIris() {
	// server
	app := iris.New()
	app.Logger().SetLevel("warn")
	app.Use(recover.New())
	app.Use(logger.New())
	app.AllowMethods(iris.MethodOptions)

	app.OnAnyErrorCode(func(ctx iris.Context) {
		path := ctx.Path()
		var err error
		if strings.Contains(path, "/api/admin/") {
			_, err = ctx.JSON(utils.JSONError(utils.RespCodeErr, "发生错误"))
		}
		if err != nil {
			logrus.Error(err)
		}
	})

	// app.Any("/", func(i iris.Context) {
	// 	_, _ = i.HTML("<h1>Powered by qwquiver</h1>")
	// })
	app.HandleDir("/", bindata.AssetFile())

	// api
	mvc.Configure(app.Party("/api"), func(m *mvc.Application) {
		m.Party("/conf").Handle(new(api.ConfController))
		m.Party("/query").Handle(new(api.QueryController))
		m.Party("/school").Handle(new(api.SchoolController))
		m.Party("/chart").Handle(new(api.ChartController))
	})

	// admin
	// mvc.Configure(app.Party("/api/admin"), func(m *mvc.Application) {
	// 	m.Router.Use(middleware.AdminAuth)
	// 	m.Party("/common").Handle(new(admin.CommonController))
	// 	m.Party("/user").Handle(new(admin.UserController))
	// })

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Instance.Port)}
	handleSignal(server)
	err := app.Run(iris.Server(server), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	}))
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logrus.Infof("got signal [%s], exiting now", s)
		if err := server.Close(); nil != err {
			logrus.Errorf("server close failed: " + err.Error())
		}

		lib.CloseDb()

		logrus.Infof("Exited")
		os.Exit(0)
	}()
}
