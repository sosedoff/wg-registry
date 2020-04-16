package service

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func fetchWireGuardPeers(iface string) ([]wgtypes.Peer, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	device, err := client.Device(iface)
	if err != nil {
		return nil, err
	}

	return device.Peers, nil
}
