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
	"fmt"
)

const (
	RegexExtractPath      = "(/(?P<zone>[^/]+))?/v2/service_instances/(?P<instance_id>[^/]+)(?P<rest>/.+)?"
	IdentityHeader        = "X-Broker-API-Originating-Identity"
	RegexProvisioningCall = "/v2/service_instances/[^/]*$"
)

type PathExtract struct {
	InstanceID string
	Zone       string
	Rest       string
}

func (p PathExtract) GeneratePath(instanceId string) string {
	return fmt.Sprintf("/v2/service_instances/%s%s", instanceId, p.Rest)
}

type HandlerConfig struct {
	ClientID           string     `mapstructure:"client_id" json:"client_id" yaml:"client_id"`
	ClientSecret       string     `mapstructure:"client_secret" json:"client_secret" yaml:"client_secret"`
	TokenURL           string     `mapstructure:"token_url" json:"token_url" yaml:"token_url"`
	ReconciliatorURL   string     `mapstructure:"reconciliator_url" json:"reconciliator_url" yaml:"reconciliator_url"`
	Scopes             []string   `mapstructure:"scopes" json:"scopes" yaml:"scopes"`
	EndpointParams     url.Values `mapstructure:"endpoint_params" json:"endpoint_params" yaml:"endpoint_params"`
	SkipVerification   bool       `mapstructure:"skip_verification" json:"skip_verification" yaml:"skip_verification"`
	UseSharedByDefault bool       `mapstructure:"use_shared_by_default" json:"use_shared_by_default" yaml:"use_shared_by_default"`
}

type ReconciliatorHandler struct {
	client             *core.ReconciliatorClient
	next               http.Handler
	useSharedByDefault bool
}

type ProvisionningParams struct {
	UseShared *bool `json:"use_shared,omitempty"`
}

func NewReconciliatorHandler(config HandlerConfig, next http.Handler) *ReconciliatorHandler {
	return &ReconciliatorHandler{core.NewReconciliatorClient(core.ReconciliatorConfig{
		ReconciliatorURL: config.ReconciliatorURL,
		ClientID:         config.ClientID,
		ClientSecret:     config.ClientSecret,
		TokenURL:         config.TokenURL,
		Scopes:           config.Scopes,
		SkipVerification: config.SkipVerification,
		EndpointParams:   config.EndpointParams,
	}), next, config.UseSharedByDefault}
}

func (h ReconciliatorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := gobis.Path(req)
	if !h.IsServiceInstanceVerb(path) {
		h.next.ServeHTTP(w, req)
		return
	}
	pathExtract := ParsePath(path)
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	sbRequest, err := ExtractSBRequest(rawBody)
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
	if ctxUser.UserName != "" {
		gobis.SetUsername(req, ctxUser.UserName)
	} else {
		gobis.SetUsername(req, ctxUser.UserID)
	}

	gobis.SetGroups(req, ctxUser.Groups...)

	recResp, found, err := h.client.Reconciliate(core.ReconciliatorRequest{
		Context: core.ContextPlatform{
			ContextUser: ctxUser,
			Context:     ctx,
		},
		InstanceID: pathExtract.InstanceID,
		Zone:       pathExtract.Zone,
		PlanID:     sbRequest.PlanID,
	})
	if err != nil {
		panic(err)
	}
	if !found {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		h.next.ServeHTTP(w, req)
		return
	}

	newPath := pathExtract.GeneratePath(recResp.InstanceID)
	r := regexp.MustCompile(RegexProvisioningCall)
	// we giving back all requests which is not a service creation/deletion (e.g.: service bindings)
	if !r.MatchString(path) || (req.Method != "PUT" && req.Method != "DELETE") {
		gobis.SetPath(req, newPath)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		h.next.ServeHTTP(w, req)
		return
	}

	err = nil
	if req.Method == "PUT" {
		h.provisioningHandler(w, req, newPath)
		return
	}
	gobis.SetPath(req, newPath)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	// we are now only on deprovisioning instance

	// We can delete safely because service has no bindings
	if recResp.BindingNB == 0 {
		h.next.ServeHTTP(w, req)
		return
	}

	// hijack when has bindings
	w.WriteHeader(http.StatusOK)

}

func (h ReconciliatorHandler) provisioningHandler(w http.ResponseWriter, req *http.Request, overridePath string) {
	brokerParams := struct {
		Parameters ProvisionningParams `json:"parameters"`
	}{}
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(rawBody, &brokerParams)
	if err != nil {
		panic(err)
	}

	useShared := brokerParams.Parameters.UseShared
	if (useShared == nil && !h.useSharedByDefault) || (useShared != nil && !(*useShared)) {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		h.next.ServeHTTP(w, req)
		return
	}
	gobis.SetPath(req, overridePath)
	var allRequest map[string]interface{}
	err = json.Unmarshal(rawBody, &allRequest)
	if err != nil {
		panic(err)
	}

	allParams := allRequest["parameters"].(map[string]interface{})
	delete(allParams, "use_shared")

	allRequest["parameters"] = allParams

	b, _ := json.Marshal(allRequest)
	req.ContentLength = int64(len(b))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	h.next.ServeHTTP(w, req)
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
		Rest:       paramsMap["rest"],
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
