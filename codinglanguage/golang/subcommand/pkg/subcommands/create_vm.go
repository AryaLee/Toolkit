package subcommands

import (
	"flag"
	"fmt"
	"os/exec"
)

var vpcBridge string

func init() {
	fs := NewCreateVmFlagSet()
	RegisterSubCommand("CreateVm", fs, CreateVm)
}

func NewCreateVmFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("CreateVm", flag.ExitOnError)
	fs.StringVar(&vpcBridge, "vpc-bridge", "vpcbr", "vpc bridge name")
	return fs
}

func CreateVm() {
	fmt.Println("run CreateVm func")
	cmd := exec.Command("ovs-vsctl", "add-port", vpcBridge, "port-name")
	if err := cmd.Run(); err != nil {
		fmt.Println("run CreateVm error: ", err)
		return
	}
}
