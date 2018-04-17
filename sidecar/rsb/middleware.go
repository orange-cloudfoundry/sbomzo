package rsb

import (
	"github.com/orange-cloudfoundry/gobis"
	"net/http"
)

type SBConfig struct {
}

type sbMiddleware struct {
}

func NewSBMiddleware() *sbMiddleware {
	return &sbMiddleware{}
}

func (m sbMiddleware) Handler(proxyRoute gobis.ProxyRoute, params interface{}, next http.Handler) (http.Handler, error) {

	return nil, nil
}

func (sbMiddleware) Schema() interface{} {
	return SBConfig{}
}
