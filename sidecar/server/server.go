package server

import (
	"github.com/orange-cloudfoundry/gobis-server/server"
	"github.com/orange-cloudfoundry/gobis"
	"net/url"
	"github.com/orange-cloudfoundry/gobis-middlewares/basicauth"
	"github.com/orange-cloudfoundry/sbomzo/sidecar/rsb"
)

func CreateServer(config Config) (*server.GobisServer, error) {
	serverConfig := &server.GobisServerConfig{
		Host:               config.Server.Host,
		Port:               config.Server.Port,
		Cert:               config.Server.Cert,
		Key:                config.Server.Key,
		LetsEncryptDomains: config.Server.LetsEncryptDomains,
		LogJson:            config.Log.InJson,
		LogLevel:           config.Log.Level,
		NoColor:            config.Log.NoColor,
	}

	builder := gobis.Builder()

	sbUrl, err := url.Parse(config.SBSource.URL)
	if err != nil {
		return nil, err
	}
	sbUrl.User = url.UserPassword(config.SBSource.Username, config.SBSource.Password)
	builder.
		AddRoute("/**", sbUrl.String()).
		WithMiddlewareParams(
		basicauth.BasicAuthConfig{
			BasicAuth: basicauth.BasicAuthOptions{
				{
					User:     config.SBSource.Username,
					Password: config.SBSource.Password,
				},
			},
		}).
		WithMiddlewareParams(
		rsb.SBConfig{
			SBOptions: &rsb.SBOptions{
				Enabled: true,
				HandlerConfig: rsb.HandlerConfig{
					SkipVerification:   config.SkipSSLVerification,
					Scopes:             config.Oauth2.Scopes,
					TokenURL:           config.Oauth2.TokenURI,
					ClientSecret:       config.Oauth2.ClientSecret,
					ReconciliatorURL:   config.Reconciliator.Endpoint,
					ClientID:           config.Oauth2.ClientId,
					UseSharedByDefault: config.Reconciliator.UseSharedByDefault,
				},
			},
		}).WithMiddlewareParams(config.MiddlewareParams)
	if config.SkipSSLVerification {
		builder.WithInsecureSkipVerify()
	}
	serverConfig.Routes = builder.Build()

	return server.NewGobisServer(serverConfig)
}
