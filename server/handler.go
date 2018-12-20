package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/infracloudio/firecracker-marathon/runtime"
)

type gatewayAPI struct {
}

func newGatewayAPI() *gatewayAPI {
	return &gatewayAPI{}
}

func (g *gatewayAPI) execute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name, ok := vars["name"]

	if !ok || name == "" {
		return
	}

	ex := runtime.NewExecutor()
	ex.CreateEnv("dummy")
	w.Write([]byte("test"))
}

func (g *gatewayAPI) uploadFunction(w http.ResponseWriter, r *http.Request) {
	uuid := "generated-uuid"
	w.Write([]byte(uuid))
}
