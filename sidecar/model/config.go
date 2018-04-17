package model

type Config struct {
	Oauth2        Oauth2        `json:"oauth2" yml:"oauth2" cloud:"oauth2"`
	SBSource      SBSource      `json:"sb_source" yml:"sb_source" cloud:"sb_source"`
	Reconciliator Reconciliator `json:"reconciliator" yml:"reconciliator" cloud:"reconciliator"`
	BasicAuth     BasicAuth     `json:"basic_auth" yml:"basic_auth" cloud:"basic_auth"`
}

type Reconciliator struct {
	Endpoint string `json:"endpoint" yml:"endpoint" cloud:"endpoint"`
}

type Oauth2 struct {
	TokenURI     string   `json:"token_uri" yml:"token_uri" cloud:"token_uri"`
	ClientId     string   `json:"client_id" yml:"client_id" cloud:"client_id"`
	ClientSecret string   `json:"client_secret" yml:"client_secret" cloud:"client_secret"`
	Scopes       []string `json:"scopes" yml:"scopes" cloud:"scopes"`
}

type SBSource struct {
	Username string `json:"username" yml:"username" cloud:"username"`
	Password string `json:"password" yml:"password" cloud:"password"`
	URL      string `json:"url" yml:"url" cloud:"url"`
}

type BasicAuth struct {
	Username string `json:"username" yml:"username" cloud:"username"`
	Password string `json:"password" yml:"password" cloud:"password"`
}
