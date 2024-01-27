package vm

import (
	"context"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	fcmodel "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/runvelocity/windhoek/internal/network"
	"github.com/runvelocity/windhoek/models"
)

type VmManager struct {
}

func (m VmManager) CreateVm(vmRequest models.FirecrackerVmRequest) (*firecracker.Machine, context.Context, error) {
	ctx := context.Background()

	cfg := m.GetConfig(vmRequest)

	// Check if kernel image is readable
	f, err := os.Open(cfg.KernelImagePath)
	if err != nil {
		return nil, ctx, err
	}
	defer f.Close()

	cmd := firecracker.VMCommandBuilder{}.WithSocketPath(cfg.SocketPath).WithBin("/usr/local/bin/firecracker").Build(ctx)

	vm, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		return nil, ctx, err
	}

	if err := vm.Start(ctx); err != nil {
		return nil, ctx, err
	}

	return vm, ctx, nil
}

func (m VmManager) GetConfig(vmRequest models.FirecrackerVmRequest) firecracker.Config {
	drives := firecracker.NewDrivesBuilder(RUNTIMES[vmRequest.Runtime]).
		// drives := firecracker.NewDrivesBuilder("/root/fcrootfs/rootfs.ext4").
		Build()
	networkInterface := firecracker.NetworkInterface{
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: network.FC_NETWORK_NAME,
			IfName:      network.FC_IF_NAME,
		},
	}
	cfg := firecracker.Config{
		SocketPath:      vmRequest.SocketPath,
		KernelImagePath: KERNEL_IMAGE_PATH,
		KernelArgs:      KERNEL_ARGS,
		Drives:          drives,
		MachineCfg: fcmodel.MachineConfiguration{
			VcpuCount:  firecracker.Int64(DEFAULT_CPU_COUNT),
			MemSizeMib: firecracker.Int64(DEFAULT_MEMORY_COUNT),
		},
		NetworkInterfaces: []firecracker.NetworkInterface{
			networkInterface,
		},
	}

	return cfg
}
