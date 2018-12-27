package runtime

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/infracloudio/firecracker-marathon/pkg/logging"
	"github.com/songgao/water"
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

	// Prepare TAP device fo runtime host
	CreateTAPDevice()
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

func CreateTAPDevice() {

	// check if TAP device currently is created or not
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "firecracker0"
	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TAP device - ", ifce.Name, " created.")

	//
	cmd := exec.Command("ip", "link", "set", "dev", "firecracker0", "up")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(cmdOutput.Bytes()))

	// fmt.Println("Does this happen ???")
	// cmd = exec.Command("ip", "addr", "add", "10.1.0.10/24", "dev", "firecracker0")
	// cmdOutput = &bytes.Buffer{}
	// cmd.Stdout = cmdOutput
	// err = cmd.Run()
	// if err != nil {
	// 	os.Stderr.WriteString(err.Error())
	// }

	fmt.Println("Does this happen again ???")
	cmd = exec.Command("brctl", "addif", "docker0", "firecracker0")
	cmdOutput = &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	fmt.Println(string(cmdOutput.Bytes()))
}
