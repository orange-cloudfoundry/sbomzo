package core

import (
	"net/url"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"crypto/tls"
	"net"
	"time"
	"golang.org/x/oauth2"
	"context"
	"strings"
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

const (
	PathReconciliate = "/v1/reconciliate"
)

type ReconciliatorConfig struct {
	ClientID         string
	ClientSecret     string
	TokenURL         string
	ReconciliatorURL string
	Scopes           []string
	EndpointParams   url.Values
	SkipVerification bool
}
type ReconciliatorClient struct {
	client           *http.Client
	reconciliatorURL string
}

func NewReconciliatorClient(config ReconciliatorConfig) *ReconciliatorClient {
	cCreds := &clientcredentials.Config{
		ClientID:       config.ClientID,
		ClientSecret:   config.ClientSecret,
		EndpointParams: config.EndpointParams,
		Scopes:         config.Scopes,
		TokenURL:       config.TokenURL,
	}
	srcTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.SkipVerification,
		},
	}
	httpClient := &http.Client{Transport: srcTransport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	authTransport := &oauth2.Transport{
		Base:   srcTransport,
		Source: cCreds.TokenSource(ctx),
	}

	return &ReconciliatorClient{&http.Client{
		Transport: authTransport,
	}, strings.TrimSuffix(config.ReconciliatorURL, "/")}
}

func (c ReconciliatorClient) Reconciliate(r ReconciliatorRequest) (recResp ReconciliateResponse, found bool, err error) {

	b, err := json.Marshal(r)
	if err != nil {
		return ReconciliateResponse{}, false, err
	}

	path := fmt.Sprintf("%s/%s", c.reconciliatorURL, PathReconciliate)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(b))
	if err != nil {
		return ReconciliateResponse{}, false, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return ReconciliateResponse{}, false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return ReconciliateResponse{
			InstanceID: r.InstanceID,
		}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return ReconciliateResponse{}, false, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReconciliateResponse{}, false, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(b, &recResp)
	if err != nil {
		return ReconciliateResponse{}, false, err
	}

	return recResp, true, err
}
