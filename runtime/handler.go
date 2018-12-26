package runtime

import (
	"net/http"

	"github.com/infracloudio/firecracker-marathon/config"
)

type runtimeAPI struct {
}

func newRuntimeAPI() *runtimeAPI {
	return &runtimeAPI{}
}

func (h *runtimeAPI) startInstance(w http.ResponseWriter, r *http.Request) {

	var cfg config.Runtime
	// if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
	// 	http.Error(w, strings.TrimSpace(err.Error()), http.StatusBadRequest)
	// 	return
	// }

	cfg = config.Runtime{
		Environment: config.Environment{
			Language: config.Go,
		},
		FunctionUUID: "xyz",
	}
	ex := NewExecutor()
	ex.StartInstance(cfg)

	w.Write([]byte("test"))
}

func (g *runtimeAPI) getInstance(w http.ResponseWriter, r *http.Request) {
	uuid := "generated-uuid"
	w.Write([]byte(uuid))
}
