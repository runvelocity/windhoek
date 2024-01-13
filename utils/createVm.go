package utils

import (
	"context"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

func CreateVm(vmRequest FirecrackerVmRequest) (*firecracker.Machine, context.Context, error) {
	ctx := context.Background()

	cfg := GetVmConfig(vmRequest)

	// Check if kernel image is readable
	f, err := os.Open(cfg.KernelImagePath)
	if err != nil {
		return nil, ctx, err
	}
	defer f.Close()

	cmd := firecracker.VMCommandBuilder{}.WithSocketPath(cfg.SocketPath).WithBin("/usr/local/bin/firecracker").Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		return nil, ctx, err
	}

	if err := m.Start(ctx); err != nil {
		return nil, ctx, err
	}

	return m, ctx, nil
}
