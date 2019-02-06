package drive

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	SecondaryDiskPath      = "/tmp/firecracker-disk-"
	MountPathSecondaryDisk = "/tmp/mount"
)

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

func CreateSecondaryDrive(uuid string) string {

	// This method initializes a secondary drive to be used
	// for attaching to firecracker

	// 	- dd if=/dev/zero of=/path/to/your/drive bs=1M count=$SIZE_IN_MB
	// - mkfs.ext4 /path/to/your/drive
	// - mount /path/to/your/drive /temp/mount
	// - copy the file to this mount
	// - curl --unix /tmp/firecracker.socket -i -X PUT "http://localhost/drives/testdrive" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"drive_id\": \"testdrive\", \"path_on_host\": \"/path/to/your/drive\", \"is_root_device\": false, \"is_read_only\": false}"
	// - Start the vm
	// - mount /dev/vdb on a mount point

	// check if the secondary drive is present or not for the function.
	// if not present , then create.
	// else return the response as present and available to mount

	path := SecondaryDiskPath + uuid
	ok, err := exists(path)
	if err != nil {
		panic("Some error when checking for the function directory")
	}

	if !ok {
		// make the required directory
		// err := os.MkdirAll(path, os.ModePerm)
		// if err != nil {
		// 	panic("Some error when creating function directory")
		// }

		cmd := exec.Command("dd", "if=/dev/zero", "of="+path, "bs=1M", "count=100")
		err = cmd.Run()

		if err != nil {
			fmt.Println("Error in dd command", err)
		}

		fmt.Println("dd ran successfully .... ", err)

		cmd = exec.Command("mkfs.ext4", path)
		err = cmd.Run()

		fmt.Println("mkfs ran successfully .... ", err)

		if err != nil {
			fmt.Println("Error in mkfs.ext4 command", err)
		}
	}
	return path
}

func MountSecondaryDisk(uuid string) string {

	diskPath := SecondaryDiskPath + uuid
	mountPath := MountPathSecondaryDisk + uuid
	cmd := exec.Command("mount", diskPath, mountPath)

	fmt.Println("mountpath - ", mountPath)
	fmt.Println("Seonc - -", diskPath)
	// if err := syscall.Mount(diskPath, mountPath, "ext4", syscall.MS_BIND, ""); err != nil {
	// 	log.Fatalf("Could not create Network namespace: %s", err)
	// }

	fmt.Println("diskpath --- ", diskPath)
	fmt.Println("mountPath --- ", mountPath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in mount command", err)
		log.Fatal(err)
	}

	return mountPath
}

func UnmountSecondaryDisk(mountPath string) {

	if err := syscall.Unmount(mountPath, syscall.MNT_DETACH); err != nil {
		log.Fatalf("Could not Unmount new Network namespace: %s", err)
	}

	// cmd := exec.Command("umount", mountPath)
	// err := cmd.Run()

	// if err != nil {
	// 	fmt.Println("Error in umount command", err)
	// }
}

func getMount(mountdir string) (string, error) {
	res, err := exec.Command("sh", "-c", fmt.Sprintf("mount | grep -w %s", mountdir)).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(res[:]), "\n")
	if len(lines) != 1 {
		return "", fmt.Errorf("bad mount output")
	}
	fields := strings.Fields(lines[0])
	if len(fields) != 6 {
		return "", fmt.Errorf("bad mount line formating")
	}
	device := fields[0]
	return device, nil
}

func AttachCodeToSecondaryDisk() {

}

func CopyToSecondaryDisk(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {

		fmt.Println("error 1 - ", err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		fmt.Println("error 2 - ", err)
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		fmt.Println("error 3 - ", err)

		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		fmt.Println("error 4 - ", err)

		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	fmt.Println("error 5 - ", err)
	fmt.Println("nbytes - ", nBytes)

	return nBytes, err
}
