package network

import (
	"bytes"
	"fmt"
	"github.com/songgao/water"
	"log"
	"os"
	"os/exec"
)

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

	cmd = exec.Command("brctl", "addif", "docker0", "firecracker0")
	cmdOutput = &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	fmt.Println(string(cmdOutput.Bytes()))
}
