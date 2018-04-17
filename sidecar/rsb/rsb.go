package rsb

import (
	"net/http"
	"github.com/pivotal-cf/brokerapi"
	"context"
)

type ReverseSB struct {
	sourceSB http.Handler
}

func NewReverseSB(sourceSB http.Handler) *ReverseSB {
	return &ReverseSB{}
}

func (sb *ReverseSB) Services(w http.ResponseWriter, r *http.Request) {

}

func (sb *ReverseSB) Provision(w http.ResponseWriter, r *http.Request) {

}

func (sb *ReverseSB) Deprovision(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (sb *ReverseSB) Bind(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (sb *ReverseSB) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	panic("implement me")
}

func (sb *ReverseSB) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	panic("implement me")
}

func (sb *ReverseSB) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	panic("implement me")
}
