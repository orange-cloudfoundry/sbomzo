package main

import (
	"github.com/orange-cloudfoundry/gobis"
	"github.com/orange-cloudfoundry/gobis-middlewares/cors"
	"github.com/orange-cloudfoundry/gobis-middlewares/secure"
	"github.com/orange-cloudfoundry/gobis-middlewares/basicauth"
	"github.com/orange-cloudfoundry/gobis-middlewares/oauth2"
	"github.com/orange-cloudfoundry/gobis-middlewares/jwt"
	"github.com/orange-cloudfoundry/gobis-middlewares/casbin"
	"github.com/orange-cloudfoundry/gobis-middlewares/cbreaker"
	"github.com/orange-cloudfoundry/gobis-middlewares/ratelimit"
	"github.com/orange-cloudfoundry/gobis-middlewares/connlimit"
	"github.com/orange-cloudfoundry/gobis-middlewares/trace"
	"github.com/orange-cloudfoundry/sbomzo/sidecar/rsb"
	gobiserver "github.com/orange-cloudfoundry/gobis-server/server"
	"github.com/orange-cloudfoundry/sbomzo/sidecar/server"
	"os"
)

func init() {
	midHandlers := []gobis.MiddlewareHandler{
		cors.NewCors(),
		secure.NewSecure(),
		basicauth.NewBasicAuth(),
		oauth2.NewOauth2(),
		rsb.NewSBMiddleware(),
		jwt.NewJwt(),
		casbin.NewCasbin(),
		cbreaker.NewCircuitBreaker(),
		ratelimit.NewRateLimit(),
		connlimit.NewConnLimit(),
		trace.NewTrace(),
	}
	gobiserver.AddMiddlewareHandlers(midHandlers...)
}

func main() {
	server.NewApp().Run(os.Args)
}
