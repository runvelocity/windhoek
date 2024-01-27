package network

import (
	"fmt"
	"os"
)

const (
	FC_NETWORK_NAME = "fcnet"
	FC_IF_NAME      = "veth0"
	NETWORK_MASK    = "/24"
	SUBNET          = "192.168.127.0" + NETWORK_MASK
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
