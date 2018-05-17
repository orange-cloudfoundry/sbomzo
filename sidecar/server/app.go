package server

import (
	"fmt"
	"github.com/cloudfoundry-community/gautocloud"
	"github.com/cloudfoundry-community/gautocloud/cloudenv"
	"github.com/cloudfoundry-community/gautocloud/connectors/generic"
	"github.com/cloudfoundry-community/gautocloud/interceptor/configfile"
	"github.com/cloudfoundry-community/gautocloud/loader"
	"github.com/urfave/cli"
	"os"
	"github.com/cloudfoundry-community/gautocloud/interceptor"
)

var confFileIntercept *configfile.ConfigFileInterceptor

func init() {
	confFileIntercept = configfile.NewConfigFile()
	gautocloud.RegisterConnector(generic.NewConfigGenericConnector(
		Config{},
		confFileIntercept,
		interceptor.NewOverwrite(),
	))
}

type RSBServerApp struct {
	*cli.App
}

func NewApp() *RSBServerApp {
	app := &RSBServerApp{cli.NewApp()}
	app.Name = "sbomzo"
	app.Version = "1.0.0"
	app.Usage = "Run a sbomzo sidecar on your service-broker"
	app.ErrWriter = os.Stderr
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config-path, c",
			Value:  cloudenv.DEFAULT_CONFIG_PATH,
			Usage:  "Path to the config file (This file will not be used in a cloud env like Cloud Foundry, Heroku or kubernetes)",
			EnvVar: cloudenv.LOCAL_CONFIG_ENV_KEY,
		},
		cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "Port to listen",
		},
		cli.StringFlag{
			Name:  "cert",
			Value: "server.crt",
			Usage: "Path to a cert file or a cert content to enable https server",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "server.key",
			Usage: "Path to a key file or a key content to enable https server",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Usage: "LogConfig level to use",
		},
		cli.BoolFlag{
			Name:  "log-json, j",
			Usage: "Write log in json",
		},
		cli.BoolFlag{
			Name:  "no-color",
			Usage: "Logger will not display colors",
		},
		cli.StringSliceFlag{
			Name:  "lets-encrypt-domains, led",
			Usage: "If set server will use a certificate generated with let's encrypt, value should be your domain(s). Host and port will be overwritten to use 0.0.0.0:443",
		},
	}
	return app
}

func (a *RSBServerApp) Run(arguments []string) (err error) {
	a.Action = a.RunServer
	return a.App.Run(arguments)
}

func (a *RSBServerApp) RunServer(c *cli.Context) error {

	confPath := c.GlobalString("config-path")
	confFileIntercept.SetConfigPath(confPath)

	config := Config{
		Server: ServerConfig{
			Port:               c.GlobalInt("port"),
			Cert:               c.GlobalString("cert"),
			Key:                c.GlobalString("key"),
			LetsEncryptDomains: c.GlobalStringSlice("lets-encrypt-domains"),
		},
		Log: LogConfig{
			NoColor: c.GlobalBool("no-color"),
			Level:   c.GlobalString("log-level"),
			InJson:  c.GlobalBool("log-json"),
		},
	}

	err := gautocloud.Inject(&config)
	if err != nil {
		if _, ok := err.(loader.ErrGiveService); ok {
			return fmt.Errorf("configuration cannot be found")
		}
		return err
	}

	gobisServer, err := CreateServer(config)
	if err != nil {
		return err
	}

	return gobisServer.Run()
}
