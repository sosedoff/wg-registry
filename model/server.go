package model

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	modeSplit = "split"
	modeFull  = "full"
)

// Server holds the server configuration
type Server struct {
	ID         int       `json:"id"`
	Mode       string    `json:"mode" form:"mode"`
	Name       string    `json:"name" form:"name"`
	Interface  string    `json:"interface" form:"interface"`
	Endpoint   string    `json:"endpoint" form:"endpoint"`
	PrivateKey string    `json:"private_key"`
	PublicKey  string    `json:"public_key"`
	IPV4Net    string    `json:"ipv4_net" form:"ipv4_net"`
	IPV4Addr   string    `json:"ipv4_addr"`
	IPV6Addr   string    `json:"ipv6_addr"`
	IPV6Net    string    `json:"ipv6_net"`
	DNS        string    `json:"dns" form:"dns"`
	ListenPort int       `json:"listen_port" form:"listen_port"`
	PostUp     string    `json:"postup" form:"post_up"`
	PostDown   string    `json:"postdown" form:"post_down"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AssignPrivateKey assigns private and public keys
func (s *Server) AssignPrivateKey() error {
	if s.PrivateKey != "" {
		return nil
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}

	s.PrivateKey = key.String()
	s.PublicKey = key.PublicKey().String()

	return nil
}

// PublicAddr returns the server host:port combo
func (s *Server) PublicAddr() string {
	return fmt.Sprintf("%s:%d", s.Endpoint, s.ListenPort)
}

// IPV4Count returns the total number of IPs in the range
func (s *Server) IPV4Count() uint64 {
	_, ipv4net, err := net.ParseCIDR(s.IPV4Net)
	if err != nil {
		return 0
	}
	return cidr.AddressCount(ipv4net)
}

// IPV4Range returns first and last IP in the range
func (s *Server) IPV4Range() (net.IP, net.IP, error) {
	_, ipv4net, err := net.ParseCIDR(s.IPV4Net)
	if err != nil {
		return nil, nil, err
	}
	first, last := cidr.AddressRange(ipv4net)
	return first, last, nil
}

// Validate validates the server
func (s *Server) Validate() error {
	if s.Name == "" {
		return errors.New("Name is required")
	}
	if s.Interface == "" {
		return errors.New("Interface is required")
	}
	if s.IPV4Addr == "" {
		return errors.New("Network is required")
	}
	if s.IPV4Count() == 0 {
		return errors.New("Network CIDR is invalid")
	}
	if s.IPV4Count() < 8 {
		return errors.New("Network CIDR is too small")
	}
	return nil
}

// ServerWithDefaults returns a server with default configuration
func ServerWithDefaults() *Server {
	iface := "wg0"
	postUp := fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE; ip6tables -A FORWARD -i %s -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE", iface, iface)
	postDown := fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE; ip6tables -D FORWARD -i %s -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE", iface, iface)

	return &Server{
		Name:       "private",
		Interface:  iface,
		IPV4Net:    "10.10.0.0/24",
		IPV4Addr:   "10.10.0.0/24",
		ListenPort: 51820,
		PostUp:     postUp,
		PostDown:   postDown,
	}
}
