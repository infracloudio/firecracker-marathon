package gateway

import (
	"fmt"
	"net/http"

	"github.com/infracloudio/firecracker-marathon/pkg/config"
	"github.com/infracloudio/firecracker-marathon/pkg/gateway/client"
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
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (g *gatewayAPI) updateFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (g *gatewayAPI) deleteFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (g *gatewayAPI) getFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (g *gatewayAPI) executeFunction(w http.ResponseWriter, r *http.Request) {

	// Call runtime , providing the configuration to run.
	cfg := config.Runtime{
		Environment: config.Environment{
			Language: config.Go,
		},
		FunctionUUID: "uuid",
	}

	fmt.Println(cfg)
	// Create a client call to Runtime
	c, err := client.NewClient(defaultHost, apiVersion, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req := c.Get().Resource(getDefaultURL() + "/instance/start").Body(cfg)
	resp := req.Do()
	if resp.Error() != nil {
		http.Error(w, resp.FormatError().Error(), http.StatusInternalServerError)
	}

	fmt.Println("Whats the response ---- ", resp.Body)
}

func getDefaultURL() string {
	return "http://" + defaultHost + ":" + defaultPort
}
