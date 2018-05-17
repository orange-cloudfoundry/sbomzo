package rsb

import (
	"github.com/orange-cloudfoundry/gobis"
	"net/http"
)

type SBConfig struct {
	SBOptions *SBOptions `mapstructure:"service_broker" json:"service_broker" yaml:"service_broker"`
}

type SBOptions struct {
	HandlerConfig
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
}

type sbMiddleware struct {
}

func NewSBMiddleware() *sbMiddleware {
	return &sbMiddleware{}
}

func (m sbMiddleware) Handler(proxyRoute gobis.ProxyRoute, params interface{}, next http.Handler) (http.Handler, error) {
	config := params.(SBConfig)
	options := config.SBOptions
	if options == nil || !options.Enabled {
		return next, nil
	}

	handler := NewReconciliatorHandler(options.HandlerConfig, next)
	return handler, nil
}

func (sbMiddleware) Schema() interface{} {
	return SBConfig{}
}
