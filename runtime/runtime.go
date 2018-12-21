package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/infracloudio/firecracker-marathon/logging"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	firecrackerBinaryOverrideEnv = "FC_TEST_BIN"
	// TODO: Get rid of these sockets
	firecrackerSocketPath = "/tmp/firecracker.socket"
	firecrackerGoVmLinux = "../environments/go/hello-vmlinux.bin"
	firecrackerGoExt4 = "../environments/go/hello-rootfs.ext4"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (r *Executor) CreateEnv(language string) error {

	// TODO check from pool of VMs if asked environment is running or not
	// if yes then, dont start a new environment
	socketPath := firecrackerSocketPath

	// Create Machine instance
	cfg := firecracker.Config{
		SocketPath:  socketPath,
		LogFifo: "/tmp/log.fifo",
		LogLevel: "Debug",
		MetricsFifo: "/tmp/metrics.fifo",
		KernelImagePath: firecrackerGoVmLinux,
		KernelArgs: "console=ttyS0 reboot=k panic=1 pci=off",
		RootDrive: firecracker.BlockDevice{
			HostPath: firecrackerGoExt4,
			Mode:     "rw",
		},
		CPUCount:    int64(1),
		CPUTemplate: firecracker.CPUTemplate(firecracker.CPUTemplateT2),
		HtEnabled:   false,
		MemInMiB:    int64(128),
		Debug: 		 true,
	}

	logger := log.NewEntry(logging.NewLogger())
	m, err := firecracker.NewMachine(cfg, firecracker.WithLogger(logger))
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	log.Printf("Calling machine init")

	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	exitchan, err := m.Init(ctx)
	if err != nil {
		log.Errorf("Firecracker init returned ")
		return err
	}

	go func() {
		<-exitchan
		vmmCancel()
	}()


	log.Printf("Calling start instance")

	err = m.StartInstance(vmmCtx)
	if err != nil {
		return errors.New("can't start firecracker - make sure it's in your path.")
	}

	log.Printf("Instance creation was successful")

	//}()
	go func() {
		<- vmmCtx.Done()
	}()
	return nil
}

// func (r *Runtime) ExecuteCodeInEnvironment() error {

// 	//

// 	return nil
// }



func removeSocket(socketPath string) {
	err := os.Remove(socketPath)
	if err != nil {
		fmt.Println("Error removing socket ..")
	}
}
