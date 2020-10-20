package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ca17/go-common/common"
	"github.com/ca17/go-common/conf"
	"github.com/ca17/go-common/log"
	"github.com/ca17/go-common/tpl"
)

func StartWebserver(config conf.AppConfig, appContext *AppContext, tplrender *tpl.CommonTemplate, handler ...WebHandler) error {
	webcfg := config.GetWebConfig()
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 5,
	// }))
	e.Use(ServerRecover(config.GetWebConfig().Debug))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: config.GetAppName()+" ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${error}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(webcfg.AllowOrigins, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		// AllowHeaders: []string{"Content-Type"},
		AllowCredentials: true,
	}))
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(webcfg.Secret),
		Skipper: func(c echo.Context) bool {
			if config.IsDev() {
				return true
			}
			skips := strings.Split(webcfg.AuthSkip, ",")
			if common.InSlice(c.Request().RequestURI, skips) {
				return true
			}
			skipPrefix, ok := appContext.Get(AuthSkipPrefix)
			if ok {
				for _, p := range skipPrefix.([]string) {
					if strings.HasPrefix(c.Request().RequestURI, p) {
						return true
					}
				}
			}
			return false
		},
	}))

	// Init Handlers
	webctx := NewWebContext(appContext, &config)
	group := e.Group("")
	for _, webHandler := range handler {
		webHandler.InitRouter(webctx, group)
	}
	e.Renderer = tplrender
	e.HideBanner = true
	e.Debug = webcfg.Debug
	log.Info("try start tls server")
	err := e.StartTLS(fmt.Sprintf("%s:%d", webcfg.Host, webcfg.Port), webcfg.CertFile, webcfg.KeyFile)
	if err != nil {
		log.Warningf("start tls server error %+v", err)
		log.Info("start server")
		err = e.Start(fmt.Sprintf("%s:%d", webcfg.Host, webcfg.Port))
	}
	return err
}
