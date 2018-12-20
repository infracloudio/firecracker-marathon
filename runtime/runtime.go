package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/infracloudio/firecracker-marathon/logging"
	log "github.com/sirupsen/logrus"
)

const (
	firecrackerBinaryPath        = "firecracker"
	firecrackerBinaryOverrideEnv = "FC_TEST_BIN"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (r *Executor) CreateEnv(language string) error {

	// TODO check from pool of VMs if asked environment is running or not
	// if yes then, dont start a new environment

	socketPath := "/tmp/firecracker.socket"
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(socketPath, logger, true)
	fmt.Println(firecrackerClient)

	defer removeSocket(socketPath)

	// Create Machine instance
	cfg := firecracker.Config{
		Debug:           true,
		KernelImagePath: "../environments/go/debian-vmlinux",
		RootDrive: firecracker.BlockDevice{
			HostPath: "../environments/go/debian.ext4",
			Mode:     "ro",
		},
		SocketPath:  socketPath,
		CPUTemplate: firecracker.CPUTemplate(firecracker.CPUTemplateT2),
		CPUCount:    int64(1),
		MemInMiB:    int64(128),
		HtEnabled:   false,
	}

	ctx := context.Background()
	cmd := firecracker.VMCommandBuilder{}.
		WithSocketPath(socketPath).
		WithBin(getFirecrackerBinaryPath()).
		Build(ctx)

	m, err := firecracker.NewMachine(cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		log.Fatalf("unexpectd error: %v", err)
	}

	vmmCtx, vmmCancel := context.WithTimeout(ctx, 30*time.Second)
	defer vmmCancel()
	exitchannel := make(chan error)
	go func() {
		exitCh, err := m.Init(vmmCtx)
		if err != nil {
			close(exitchannel)
			log.Fatalf("Failed to start VMM: %v", err)
		}
		exitchannel <- <-exitCh
		close(exitchannel)
	}()
	return nil
}

// func (r *Runtime) ExecuteCodeInEnvironment() error {

// 	//

// 	return nil
// }

func getFirecrackerBinaryPath() string {
	if val := os.Getenv(firecrackerBinaryOverrideEnv); val != "" {
		return val
	}
	return filepath.Join("/usr/local/bin", firecrackerBinaryPath)
}

func removeSocket(socketPath string) {
	err := os.Remove(socketPath)
	if err != nil {
		fmt.Println("Error removing socket ..")
	}
}
