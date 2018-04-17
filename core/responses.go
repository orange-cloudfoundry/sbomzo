package core

const (
	ExposurePublic  Exposure = "public"
	ExposurePrivate Exposure = "private"
)

type Exposure string

func (e Exposure) IsPublic() bool {
	return e == ExposurePublic
}

func (e Exposure) IsPrivate() bool {
	return e == ExposurePrivate
}

type ReconciliateResponse struct {
	InstanceID  string `json:"instance_id"`
	MasterZone  *Zone  `json:"master_zone,omitempty"`
	CurrentZone *Zone  `json:"current_zone,omitempty"`
}

type Zone struct {
	InstanceID string          `json:"instance_id,omitempty"`
	ServiceID  string          `json:"service_id,omitempty"`
	Platform   string          `json:"platform,omitempty"`
	Exposure   Exposure        `json:"exposure,omitempty"`
	Site       string          `json:"site,omitempty"`
	Context    ContextPlatform `json:"context,omitempty"`
}
