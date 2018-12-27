package gateway

import (
	"net/http"

	"github.com/infracloudio/firecracker-marathon/pkg/config"
	"github.com/infracloudio/firecracker-marathon/pkg/runtime"
)

const (
	defaultHost = "localhost"
	defaultPort = "8383"
	apiVersion  = "v1"
)

type gatewayAPI struct {
}

func newGatewayAPI() *gatewayAPI {
	return &gatewayAPI{}
}

func (g *gatewayAPI) uploadFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (g *gatewayAPI) updateFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (g *gatewayAPI) deleteFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (g *gatewayAPI) getFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (g *gatewayAPI) executeFunction(w http.ResponseWriter, r *http.Request) {




	cfg := config.Runtime{
		Environment: config.Environment{
			Language: config.Go,
		},
		FunctionUUID: "uuid",
	}
	ex := runtime.NewExecutor()
	ex.StartInstance(cfg)

	w.Write([]byte("test"))
}

func getDefaultURL() string {
	return "http://" + defaultHost + ":" + defaultPort
}
