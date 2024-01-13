package utils

import (
	"fmt"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

var (
	ROOTFS_PATH                            = "/root/fcrootfs"
	KERNEL_IMAGE_PATH                      = "/root/fckernels/vmlinux"
	FC_SOCKETS_PATH                        = "/root/fcsockets"
	FC_NETWORK_NAME                        = "fcnet"
	FC_IF_NAME                             = "veth0"
	KERNEL_ARGS                            = "console=ttyS0 reboot=k panic=1 pci=off"
	NETWORK_MASK                           = "/24"
	SUBNET                                 = "192.168.127.0" + NETWORK_MASK
	DEFAULT_CPU_COUNT    int64             = 1
	DEFAULT_MEMORY_COUNT int64             = 512
	AWS_REGION                             = "us-east-1"
	RUNTIMES             map[string]string = map[string]string{
		"nodejs": "/root/fcruntimes/nodejs-runtime.ext4",
	}
)

func WriteCNIConfWithHostLocalSubnet(path string) error {
	return os.WriteFile(path, []byte(fmt.Sprintf(
		`{
	"cniVersion": "1.0.0",
	"name": "%s",
	"plugins": [
		{
		"type": "ptp",
		"ipMasq": true,
		"ipam": {
			"type": "host-local",
			"subnet": "%s"
		}
		},
		{
			"type": "firewall"
		},
		{
		"type": "tc-redirect-tap"
		}
	]
	}`, FC_NETWORK_NAME, SUBNET)), 0644)
}

func GetVmConfig(vmRequest FirecrackerVmRequest) firecracker.Config {
	drives := firecracker.NewDrivesBuilder(RUNTIMES[vmRequest.Runtime]).
		// drives := firecracker.NewDrivesBuilder("/root/fcrootfs/rootfs.ext4").
		Build()
	networkInterface := firecracker.NetworkInterface{
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: FC_NETWORK_NAME,
			IfName:      FC_IF_NAME,
		},
	}
	cfg := firecracker.Config{
		SocketPath:      vmRequest.SocketPath,
		KernelImagePath: KERNEL_IMAGE_PATH,
		KernelArgs:      KERNEL_ARGS,
		Drives:          drives,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(DEFAULT_CPU_COUNT),
			MemSizeMib: firecracker.Int64(DEFAULT_MEMORY_COUNT),
		},
		NetworkInterfaces: []firecracker.NetworkInterface{
			networkInterface,
		},
	}

	return cfg
}
