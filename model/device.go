package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Device contains the device details
type Device struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	OS                  string    `json:"os"`
	Name                string    `json:"name"`
	Enabled             bool      `json:"enabled"`
	PrivateKey          string    `json:"private_key"`
	PublicKey           string    `json:"public_key"`
	IPV4                string    `json:"ipv4"`
	IPV6                string    `json:"ipv6"`
	PersistentKeepalive int       `json:"persistent_keepalive"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	peerInfo *wgtypes.Peer
}

func (d *Device) SetPeerInfo(info *wgtypes.Peer) {
	d.peerInfo = info
}

func (d *Device) GetPeerInfo() *wgtypes.Peer {
	return d.peerInfo
}

// AssignPrivateKey assigns private and public keys
func (d *Device) AssignPrivateKey() error {
	if d.PrivateKey != "" {
		return nil
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}

	d.PrivateKey = key.String()
	d.PublicKey = key.PublicKey().String()

	return nil
}

// IconClass returns an icon class based on type of device
func (d *Device) IconClass() string {
	switch d.OS {
	case "mac", "ios":
		return "fab fa-apple"
	case "windows":
		return "fab fa-windows"
	case "android":
		return "fab fa-android"
	default:
		return "fas fa-server"
	}
}

// IsMobile returns true if devices is of a mobile kind
func (d *Device) IsMobile() bool {
	return d.OS == "ios" || d.OS == "android"
}

// AllowedIPs returns a set of allowed IP addresse
func (d *Device) AllowedIPs() string {
	return fmt.Sprintf("%s/32", d.IPV4)
}

// PeerIP returns a peering IP address
func (d *Device) PeerIP(network string) string {
	netmask := strings.Split(network, "/")[1]
	return fmt.Sprintf("%s/%s", d.IPV4, netmask)
}

// Validate checks device validity
func (d *Device) Validate() error {
	if d.UserID <= 0 {
		return errors.New("User is required")
	}
	if d.Name == "" {
		return errors.New("Name is required")
	}
	if d.OS == "" {
		return errors.New("OS is required")
	}
	if d.PrivateKey == "" {
		return errors.New("Private key is required")
	}
	if d.PublicKey == "" {
		return errors.New("Public key is required")
	}
	return nil
}
