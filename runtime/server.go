package runtime

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/infracloudio/firecracker-marathon/logging"
)

func Start() {
	r := newRuntimeAPI()
	routes := r.Routes()
	startServer(routes)
}

func startServer(routes []*Route) {
	var err error
	logger := logging.NewLogger()

	router := mux.NewRouter()

	for _, v := range routes {
		router.Methods(v.verb).Path(v.path).HandlerFunc(v.fn)
	}

	server := &http.Server{
		Addr:    ":8383",
		Handler: router,
	}

	logger.Info("Runtime Server listening")
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Runtime Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatalf("Runtime Server Shutdown: %s", err)
	}
	logger.Info("Runtime Server exiting")
}
