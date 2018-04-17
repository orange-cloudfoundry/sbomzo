// Reverse proxy to another service broker
package rsb

import (
	"net/http"
	"net/url"
	"github.com/orange-cloudfoundry/sbomzo/core"
	"strings"
	"github.com/orange-cloudfoundry/gobis"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

const (
	RegexExtractPath = "(/(?P<zone>[^/]+))?/v2/[^/]+/(?P<instance_id>[^/]+)"
	IdentityHeader   = "X-Broker-API-Originating-Identity"
)

type PathExtract struct {
	InstanceID string
	Zone       string
}

type ReconciliatorConfig struct {
	ClientID         string
	ClientSecret     string
	TokenURL         string
	ReconciliatorURL string
	Scopes           []string
	EndpointParams   url.Values
	SkipVerification bool
}

type ReconciliatorHandler struct {
	client *core.ReconciliatorClient
	next   http.Handler
}

func NewReconciliatorHandler(config ReconciliatorConfig, next http.Handler) *ReconciliatorHandler {
	return &ReconciliatorHandler{core.NewReconciliatorClient(core.ReconciliatorConfig{
		ReconciliatorURL: config.ReconciliatorURL,
		ClientID:         config.ClientID,
		ClientSecret:     config.ClientSecret,
		TokenURL:         config.TokenURL,
		Scopes:           config.Scopes,
		SkipVerification: config.SkipVerification,
		EndpointParams:   config.EndpointParams,
	}), next}
}

func (h ReconciliatorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := gobis.Path(req)
	if !h.IsServiceInstanceVerb(path) {
		h.next.ServeHTTP(w, req)
		return
	}
	extractPath := ParsePath(path)
	rawReq, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	sbRequest, err := ExtractSBRequest(rawReq)
	if err != nil {
		panic(err)
	}
	var profile core.Profile
	identityRaw := req.Header.Get(IdentityHeader)
	if identityRaw != "" {
		profile, err = core.ExtractProfile(identityRaw)
		if err != nil {
			panic(err)
		}
	}

	ctx := sbRequest.Context
	if ctx.Platform == "" {
		ctx.Platform = profile.Platform
	}

	ctxUser := profile.ToContextUser()
	gobis.SetUsername(req, ctxUser.UserName)
	gobis.SetGroups(req, ctxUser.Groups...)

	recResp, found, err := h.client.Reconciliate(core.ReconciliatorRequest{
		Context: core.ContextPlatform{
			ContextUser: ctxUser,
			Context:     ctx,
		},
		InstanceID: extractPath.InstanceID,
		Zone:       extractPath.Zone,
		PlanID:     sbRequest.PlanID,
	})
	if err != nil {
		panic(err)
	}
	if !found {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawReq))
		h.next.ServeHTTP(w, req)
		return
	}

}

func (ReconciliatorHandler) IsServiceInstanceVerb(path string) bool {
	return strings.Contains(path, "/v2/service_instances")
}

func ParsePath(path string) PathExtract {
	r := regexp.MustCompile(RegexExtractPath)
	sMatch := r.FindStringSubmatch(path)

	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(sMatch) {
			paramsMap[name] = sMatch[i]
		}
	}
	return PathExtract{
		Zone:       paramsMap["zone"],
		InstanceID: paramsMap["instance_id"],
	}
}

func ExtractSBRequest(b []byte) (core.SBRequest, error) {
	var sbRequest core.SBRequest
	err := json.Unmarshal(b, &sbRequest)
	if err != nil {
		return core.SBRequest{}, err
	}
	return sbRequest, nil
}
