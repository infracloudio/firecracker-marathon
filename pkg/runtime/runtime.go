package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/infracloudio/firecracker-marathon/pkg/config"
	"github.com/infracloudio/firecracker-marathon/pkg/drive"
	"github.com/infracloudio/firecracker-marathon/pkg/logging"
	log "github.com/sirupsen/logrus"
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

func (r *Executor) StartInstance(c config.Runtime) {
	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()

	action := &models.InstanceActionInfo{
		ActionType: "InstanceStart",
	}
	a, err := firecrackerClient.CreateSyncAction(ctx, action)
	if err != nil {
		fmt.Println("Error in starting instance", err)
		fmt.Println("Another error --", a.Error())
	}

}

func (r *Executor) StopInstance(c config.Runtime) {

	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()
	defer r.RemoveSocket(conf.SocketPath)

	action := &models.InstanceActionInfo{
		ActionType: "InstanceHalt",
	}
	a, err := firecrackerClient.CreateSyncAction(ctx, action)
	if err != nil {
		fmt.Println("Error in stoppping instance", err)
		fmt.Println("Another error --", a.Error())
		log.Fatal("Error in instance halt")
	}
}

func (r *Executor) OpenSocket(c config.Runtime) {
	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())

	errCh := make(chan error)
	err := waitForSocket(conf.SocketPath, logger, 3*time.Second, errCh)
	if err != nil {
		msg := fmt.Sprintf("Firecracker did not create API socket %s: %s", conf.SocketPath, err)
		err = fmt.Errorf("Error in creating socketpath --- ", msg)
		log.Fatal(err)
	}
}

// open Socket connection
func waitForSocket(socketPath string, logger *log.Entry, timeout time.Duration, exitchan chan error) error {

	firecrackerClient := firecracker.NewFirecrackerClient(socketPath, logger, true)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error)
	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			case err := <-exitchan:
				done <- err
				return
			case <-ticker.C:
				if _, err := os.Stat(socketPath); err != nil {
					fmt.Println("what is err --- ", err)
					continue
				}

				// Send test HTTP request to make sure socket is available
				if _, err := firecrackerClient.GetMachineConfig(); err != nil {
					fmt.Println("What is error --- ", err)
					continue
				}

				done <- nil
				return
			}
		}
	}()

	return <-done
}
func (r *Executor) AttachUserCodeToSecondaryDisk(c config.Runtime) {

	//	diskPath := drive.SecondaryDiskPath + c.FunctionUUID

	// // Mount the Disk
	// mountPath := drive.MountSecondaryDisk(c.FunctionUUID)

	// // Add the code that is present in some directory , asusumption for now.
	// path := filepath.Join(UserCodeDirectory, c.FunctionUUID+".go")

	// drive.CopyToSecondaryDisk(path, mountPath)
	// // Unmount the Disk
	// drive.UnmountSecondaryDisk(mountPath)

	// conf := InitializeRunConfig(c)
	// logger := log.NewEntry(logging.NewLogger())
	// firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	// ctx := context.Background()

	// driveID := "drive-" + c.FunctionUUID
	// isReadonly := false
	// isRootDevice := false

	// drive := &models.Drive{
	// 	DriveID:      &driveID,
	// 	IsReadOnly:   &isReadonly,
	// 	IsRootDevice: &isRootDevice,
	// 	PathOnHost:   &diskPath,
	// }

	// a, err := firecrackerClient.PutGuestDriveByID(ctx, driveID, drive)
	// if err != nil {
	// 	fmt.Println("Error in attaching drive", err)
	// 	fmt.Println("Another error --", a.Error())
	// }
}

// func (r *Executor) CopyIPDetailsToSecondaryDisk(c config.Runtime) {

// 	diskPath := drive.SecondaryDiskPath + c.FunctionUUID

// 	// Mount the Disk
// 	mountPath := drive.MountSecondaryDisk(c.FunctionUUID)

// 	// Add the code that is present in some directory , asusumption for now.
// 	path := filepath.Join(MachineDetails, c.FunctionUUID)

// 	//network.ProvideIPToVM(path, mountPath)
// 	// Unmount the Disk
// 	drive.UnmountSecondaryDisk(mountPath)

// 	conf := InitializeRunConfig(c)
// 	logger := log.NewEntry(logging.NewLogger())
// 	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

// 	ctx := context.Background()

// 	driveID := "drive-" + c.FunctionUUID
// 	isReadonly := true
// 	isRootDevice := false

// 	drive := &models.Drive{
// 		DriveID:      &driveID,
// 		IsReadOnly:   &isReadonly,
// 		IsRootDevice: &isRootDevice,
// 		PathOnHost:   &diskPath,
// 	}

// 	a, err := firecrackerClient.PutGuestDriveByID(ctx, driveID, drive)
// 	if err != nil {
// 		fmt.Println("Error in attaching drive", err)
// 		fmt.Println("Another error --", a.Error())
// 	}
// }

func (r *Executor) AddBootSource(c config.Runtime) {

	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()

	source := &models.BootSource{
		BootArgs:        conf.KernelOpts,
		KernelImagePath: &conf.KernelPath,
	}
	a, err := firecrackerClient.PutGuestBootSource(ctx, source)

	if err != nil {
		fmt.Println("Error in adding bootsource", err)
		fmt.Println("Another error --", a.Error())
	}

}

func (r *Executor) AddSecondaryDrive(c config.Runtime) {
	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()

	driveID := "secondary-drive"
	isReadonly := false
	isRootDevice := false

	secondaryPath := drive.CreateSecondaryDrive(c.FunctionUUID)

	// Mount the Disk
	mountPath := drive.MountSecondaryDisk(c.FunctionUUID)

	// Add the code that is present in some directory , asusumption for now.
	path := filepath.Join(UserCodeDirectory, c.FunctionUUID+".go")

	_, err := drive.CopyToSecondaryDisk(path, mountPath+"/"+c.FunctionUUID+".go")
	if err != nil {
		log.Fatal("there is error ;;;; ", err)
	}

	//time.Sleep(3 * time.Minute)
	// Unmount the Disk
	drive.UnmountSecondaryDisk(mountPath)

	drive := &models.Drive{
		DriveID:      &driveID,
		IsReadOnly:   &isReadonly,
		IsRootDevice: &isRootDevice,
		PathOnHost:   &secondaryPath,
	}

	a, err := firecrackerClient.PutGuestDriveByID(ctx, driveID, drive)
	if err != nil {
		fmt.Println("Error in attaching root drive", err)
		fmt.Println("Another error --", a.Error())
	}
}

func (r *Executor) AddRootDrive(c config.Runtime) {
	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()

	driveID := "root-drive"
	isReadonly := true
	isRootDevice := true

	drive := &models.Drive{
		DriveID:      &driveID,
		IsReadOnly:   &isReadonly,
		IsRootDevice: &isRootDevice,
		PathOnHost:   &conf.RootDrive,
	}

	a, err := firecrackerClient.PutGuestDriveByID(ctx, driveID, drive)
	if err != nil {
		fmt.Println("Error in attaching root drive", err)
		fmt.Println("Another error --", a.Error())
	}
}

func (r *Executor) AddNetworkInterface(c config.Runtime) {

	conf := InitializeRunConfig(c)
	logger := log.NewEntry(logging.NewLogger())
	firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)

	ctx := context.Background()

	ifaceID := "eth0"
	ifaceCfg := &models.NetworkInterface{
		AllowMmdsRequests: true,
		HostDevName:       "firecracker0",
		IfaceID:           &ifaceID,
	}

	a, err := firecrackerClient.PutGuestNetworkInterfaceByID(ctx, ifaceID, ifaceCfg)
	if err != nil {
		fmt.Println("Error in attaching network interface", err)
		fmt.Println("Another error --", a.Error())
	}

}

// func (r *Executor) ProvideIPToVM(c config.Runtime) {

// }
// func (r *Executor) StartInstance(c config.Runtime) error {

// TODO check from pool of VMs if asked environment is running or not
// if yes then, dont start a new environment

// conf := InitializeRunConfig(c)
// logger := log.NewEntry(logging.NewLogger())
// firecrackerClient := firecracker.NewFirecrackerClient(conf.SocketPath, logger, true)
// fmt.Println(firecrackerClient)
// defer removeSocket(conf.SocketPath)

// diskPath := drive.SecondaryDiskPath + c.FunctionUUID

// // Create Machine instance
// cfg := firecracker.Config{
// 	Debug:           true,
// 	KernelImagePath: conf.KernelPath,

// 	RootDrive: firecracker.BlockDevice{
// 		HostPath: conf.RootDrive,
// 		Mode:     "ro",
// 	},
// 	SocketPath:  conf.SocketPath,
// 	CPUTemplate: firecracker.CPUTemplate(firecracker.CPUTemplateT2),
// 	CPUCount:    int64(1),
// 	MemInMiB:    int64(128),
// 	HtEnabled:   false,
// 	AdditionalDrives: []firecracker.BlockDevice{
// 		firecracker.BlockDevice{
// 			HostPath: diskPath,
// 			Mode:     "ro",
// 		},
// 	},
// 	NetworkInterfaces: []firecracker.NetworkInterface{
// 		&firecracker.NetworkInterface{
// 			AllowMDDS:   true,
// 			HostDevName: "firecracker0",
// 		},
// 	},
// }

// ctx := context.Background()
// m, err := firecracker.NewMachine(cfg, firecracker.WithLogger(logger))
// if err != nil {
// 	log.Errorf("unexpected error: %v", err)
// 	return err
// }

// vmmCtx, vmmCancel := context.WithCancel(ctx)
// exitchan, err := m.Init(vmmCtx)
// if err != nil {
// 	removeSocket(conf.SocketPath)
// 	fmt.Printf("Firecracker Init returned error %s", err)
// 	return err
// }

// go func() {
// 	<-exitchan
// 	removeSocket(conf.SocketPath)
// 	vmmCancel()
// }()

// // err = m.StartInstance(vmmCtx)
// // if err != nil {
// // 	fmt.Println("err --- ", err)
// // 	return errors.New("can't start firecracker - make sure it's in your path.")
// // }
// 	return nil
// }

func (r *Executor) RemoveSocket(socketPath string) {
	err := os.Remove(socketPath)
	if err != nil {
		fmt.Println("Error removing socket ..")
	}
}
