package server

type Config struct {
	Log                 LogConfig              `json:"log" yaml:"log" cloud:"log"`
	Server              ServerConfig           `json:"server" yaml:"server" cloud:"server"`
	Oauth2              Oauth2Config           `json:"oauth2" yaml:"oauth2" cloud:"oauth2"`
	SBSource            SBSourceConfig         `json:"sb_source" yaml:"sb_source" cloud:"sb_source"`
	Reconciliator       ReconciliatorConfig    `json:"reconciliator" yaml:"reconciliator" cloud:"reconciliator"`
	SkipSSLVerification bool                   `json:"skip_ssl_verification" yaml:"skip_ssl_verification" cloud:"skip_ssl_verification"`
	MiddlewareParams    map[string]interface{} `json:"middleware_params" yaml:"middleware_params"`
}

type LogConfig struct {
	Level   string `json:"level" yaml:"level"`
	InJson  bool   `json:"in_json" yaml:"in_json"`
	NoColor bool   `json:"no_color" yaml:"no_color"`
}

type ServerConfig struct {
	Host               string   `json:"host" yaml:"host"`
	Port               int      `json:"port" yaml:"port"`
	Cert               string   `json:"cert" yaml:"cert" cloud-default:"server.crt"`
	Key                string   `json:"key" yaml:"key" cloud-default:"server.key"`
	LetsEncryptDomains []string `json:"lets_encrypt_domains" yaml:"lets_encrypt_domains"`
}

type ReconciliatorConfig struct {
	Endpoint           string `json:"endpoint" yaml:"endpoint" cloud:"endpoint"`
	UseSharedByDefault bool   `json:"use_shared_by_default" yaml:"use_shared_by_default"`
}

type Oauth2Config struct {
	TokenURI     string   `json:"token_uri" yaml:"token_uri" cloud:"token_uri"`
	ClientId     string   `json:"client_id" yaml:"client_id" cloud:"client_id"`
	ClientSecret string   `json:"client_secret" yaml:"client_secret" cloud:"client_secret"`
	Scopes       []string `json:"scopes" yaml:"scopes" cloud:"scopes"`
}

type SBSourceConfig struct {
	Username string `json:"username" yaml:"username" cloud:"username"`
	Password string `json:"password" yaml:"password" cloud:"password"`
	URL      string `json:"url" yaml:"url" cloud:"url"`
}
