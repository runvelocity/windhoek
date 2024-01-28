package vm

var (
	KERNEL_IMAGE_PATH                      = "/root/fckernels/vmlinux"
	KERNEL_ARGS                            = "console=ttyS0 reboot=k panic=1 pci=off"
	DEFAULT_CPU_COUNT    int64             = 1
	DEFAULT_MEMORY_COUNT int64             = 512
	FC_SOCKETS_PATH                        = "/root/fcsockets"
	RUNTIMES             map[string]string = map[string]string{
		"nodejs": "/root/fcruntimes/nodejs-runtime.ext4",
	}
)
