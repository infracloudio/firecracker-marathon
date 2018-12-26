package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/infracloudio/firecracker-marathon/pkg/config"
	"github.com/infracloudio/firecracker-marathon/pkg/logging"
	log "github.com/sirupsen/logrus"
)

const (
	KernelAndRootFSLocation = "/var/environments"
	RootFS                  = "root-go.ext4"
	Kernel                  = "vmlinux-go.bin"
	FirecrackerBinaryPath   = "/usr/local/bin"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func InitializeRunConfig(c config.Runtime) config.RunConfig {

	kernelPath, rootFSPath := getKernelFilePath(c.Environment)
	socketPath := "/tmp/firecracker-" + c.FunctionUUID + ".socket"
	kernelOpts := "console=ttyS0 reboot=k panic=1 pci=off"

	return config.RunConfig{
		KernelOpts: kernelOpts,
		KernelPath: kernelPath,
		RootDrive:  rootFSPath,
		Language:   c.Environment.Language,
		SocketPath: socketPath,
	}
}

func getKernelFilePath(env config.Environment) (string, string) {

	var path string
	switch env.Language {
	case config.Go:
		path = filepath.Join(KernelAndRootFSLocation, "go")
	case config.Python:
		path = filepath.Join(KernelAndRootFSLocation, "python")
	default:
		return "", ""
	}

	kernelPath := filepath.Join(path, Kernel)
	ok, err := exists(filepath.Join(path, Kernel))
	if err != nil || !ok {
		fmt.Println("err or Kernel file not found", err)
		return "", ""
	}
	rootFSPath := filepath.Join(path, RootFS)
	ok, err = exists(filepath.Join(path, RootFS))
	if err != nil || !ok {
		fmt.Println("err or Root fs not found", err)
		return "", ""
	}
	return kernelPath, rootFSPath
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (r *Executor) CreateNetworkDeviceOnHost() {

}

func (r *Executor) StartInstance(c config.Runtime) error {

	// TODO check from pool of VMs if asked environment is running or not
	// if yes then, dont start a new environment

	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)
	fmt.Println(firecrackerClient)
	defer removeSocket(conf.SocketPath)

	// Create Machine instance
	cfg := firecracker.Config{
		Debug:           true,
		KernelImagePath: conf.KernelPath,
		RootDrive: firecracker.BlockDevice{
			HostPath: conf.RootDrive,
			Mode:     "ro",
		},
		SocketPath:  conf.SocketPath,
		CPUTemplate: firecracker.CPUTemplate(firecracker.CPUTemplateT2),
		CPUCount:    int64(1),
		CPUTemplate: firecracker.CPUTemplate(firecracker.CPUTemplateT2),
		HtEnabled:   false,
		MemInMiB:    int64(128),
		Debug:       true,
	}

	logger := log.NewEntry(logging.NewLogger())
	m, err := firecracker.NewMachine(cfg, firecracker.WithLogger(logger))
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	log.Printf("Calling machine init")

	ctx := context.Background()
	m, err := firecracker.NewMachine(cfg, firecracker.WithLogger(logger))
	if err != nil {
		log.Errorf("unexpected error: %v", err)
		return err
	}

	vmmCtx, vmmCancel := context.WithCancel(ctx)
	exitchan, err := m.Init(vmmCtx)
	if err != nil {
		removeSocket(conf.SocketPath)
		fmt.Printf("Firecracker Init returned error %s", err)
		return err
	}

	go func() {
		<-exitchan
		removeSocket(conf.SocketPath)
		vmmCancel()
	}()

	// err = m.StartInstance(vmmCtx)
	// if err != nil {
	// 	fmt.Println("err --- ", err)
	// 	return errors.New("can't start firecracker - make sure it's in your path.")
	// }
	return nil
}

// func (r *Runtime) ExecuteCodeInEnvironment() error {

// 	//

// 	return nil
// }

func getFirecrackerBinaryPath() string {
	return filepath.Join(FirecrackerBinaryPath, "firecracker")
}

func removeSocket(socketPath string) {
	err := os.Remove(socketPath)
	if err != nil {
		fmt.Println("Error removing socket ..")
	}
}
