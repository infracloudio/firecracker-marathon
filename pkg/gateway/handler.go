package gateway

import (
	"net/http"
	"time"

	"github.com/infracloudio/firecracker-marathon/pkg/config"
	"github.com/infracloudio/firecracker-marathon/pkg/runtime"
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

	cfg = config.Runtime{
		Environment: config.Environment{
			Language: config.Go,
		},
		FunctionUUID: "xyz",
	}

	//ex.OpenSocket(cfg)
	ex.AddBootSource(cfg)
	ex.AddRootDrive(cfg)

	ex.AddSecondaryDrive(cfg)

	ex.AttachUserCodeToSecondaryDisk(cfg)

	ex.AddNetworkInterface(cfg)

	ex.StartInstance(cfg)

	time.Sleep(2 * time.Minute)

	ex.StopInstance(cfg)

	w.Write([]byte("test"))
}
