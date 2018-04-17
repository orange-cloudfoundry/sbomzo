package core

import (
	"strings"
	"encoding/base64"
	"encoding/json"
)

const (
	PlatformCF  = "cloudfoundry"
	PlatformK8s = "kubernetes"
)

type SBRequest struct {
	ServiceID        string                 `json:"service_id"`
	PlanID           string                 `json:"plan_id"`
	Context          Context                `json:"context"`
	OrganizationGUID string                 `json:"organization_guid"`
	SpaceGUID        string                 `json:"space_guid"`
	Parameters       map[string]interface{} `json:"parameters"`
}

type ReconciliatorRequest struct {
	InstanceID string          `json:"instance_id"`
	PlanID     string          `json:"plan_id"`
	Zone       string          `json:"zone,omitempty"`
	Context    ContextPlatform `json:"context,omitempty"`
}

type Profile struct {
	Platform string `json:"platform"`
	ProfileCF
	ProfileK8S
}

type ProfileCF struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"user_name,omitempty"`
}

type ProfileK8S struct {
	Username string   `json:"username,omitempty"`
	UID      string   `json:"uid,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Extra    map[string][]string
}

func (p Profile) IsCloudFoundry() bool {
	return p.Platform == PlatformCF
}

func (p Profile) IsKubernetes() bool {
	return p.Platform == PlatformK8s
}

func (p Profile) ToContextUser() ContextUser {
	if p.IsCloudFoundry() {
		return ContextUser{
			UserName: p.ProfileCF.Username,
			UserID:   p.ProfileCF.UserID,
			Groups:   []string{},
		}
	}
	return ContextUser{
		UserName: p.ProfileK8S.Username,
		UserID:   p.ProfileK8S.UID,
		Groups:   p.ProfileK8S.Groups,
	}
}

func ExtractProfile(headerValue string) (Profile, error) {
	tokens := strings.Split(headerValue, " ")
	platform := tokens[0]
	if len(tokens) == 1 {
		return Profile{
			Platform: platform,
		}, nil
	}
	b, err := base64.RawStdEncoding.DecodeString(tokens[1])
	if err != nil {
		return Profile{}, err
	}

	var prof Profile
	err = json.Unmarshal(b, &prof)
	if err != nil {
		return Profile{}, err
	}
	prof.Platform = platform
	return prof, nil
}

type Context struct {
	ContextCF
	ContextK8S
	Platform string `json:"platform"`
}

func (c Context) IsCloudFoundry() bool {
	return c.Platform == PlatformCF
}

func (c Context) IsKubernetes() bool {
	return c.Platform == PlatformK8s
}

type ContextPlatform struct {
	ContextUser
	Context
}

type ContextUser struct {
	UserID   string   `json:"user_id,omitempty"`
	UserName string   `json:"user_name,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type ContextCF struct {
	OrganizationGUID string `json:"organization_guid,omitempty"`
	SpaceGUID        string `json:"space_guid,omitempty"`
	Endpoint         string `json:"endpoint,omitempty"`
	ServiceGUID      string `json:"service_guid,omitempty"`
}

type ContextK8S struct {
	Namespace string `json:"namespace,omitempty"`
	ClusterID string `json:"clusterid,omitempty"`
}
