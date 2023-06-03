package main

import (
	"crypto/subtle"
	"encoding/base64"
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"strings"
)

var (
	flagListenAddr   = flag.String("listen", ":33080", "the http address to start the proxy server on")
	flagAuthUser     = flag.String("auth-user", "proxy", "the basic-auth user")
	flagAuthPassword = flag.String("auth-password", "password", "the basic-auth password")
)

func main() {
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			proxyAuth := c.Request().Header.Get("Proxy-Authorization")
			if proxyAuth == "" || !strings.HasPrefix(proxyAuth, "Basic ") {
				e.Logger.Error("no Proxy-Authorization found")
				return c.NoContent(401)
			}

			authDecoded, err := base64.StdEncoding.DecodeString(strings.SplitN(proxyAuth, " ", 2)[1])
			if err != nil {
				e.Logger.Errorf("auth base64: %s", err.Error())
				return c.NoContent(401)
			}

			keyParts := strings.SplitN(string(authDecoded), ":", 2)
			if len(keyParts) != 2 {
				e.Logger.Error("invalid user or password")
				return c.NoContent(401)
			}

			user, pass := keyParts[0], keyParts[1]

			if subtle.ConstantTimeCompare([]byte(user), []byte(*flagAuthUser)) != 1 {
				e.Logger.Error("invalid user specified")
				return c.NoContent(401)
			}

			if subtle.ConstantTimeCompare([]byte(pass), []byte(*flagAuthPassword)) != 1 {
				e.Logger.Error("invalid pass specified")
				return c.NoContent(401)
			}

			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.WrapHandler(proxy)
	})

	log.Fatal(e.Start(*flagListenAddr))
}
