package model

import (
	"fmt"
	"net"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
)

type Server struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Interface  string    `json:"interface"`
	Endpoint   string    `json:"endpoint"`
	PrivateKey string    `json:"private_key"`
	PublicKey  string    `json:"public_key"`
	IPV4Net    string    `json:"ipv4_net"`
	IPV4Addr   string    `json:"ipv4_addr"`
	IPV6Addr   string    `json:"ipv6_addr"`
	IPV6Net    string    `json:"ipv6_net"`
	DNS        string    `json:"dns"`
	ListenPort int       `json:"listen_port"`
	PostUp     string    `json:"postup"`
	PostDown   string    `json:"postdown"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *Server) PublicAddr() string {
	return fmt.Sprintf("%s:%d", s.Endpoint, s.ListenPort)
}

func (s *Server) IPV4Range() (net.IP, net.IP, error) {
	_, ipv4net, err := net.ParseCIDR(s.IPV4Net)
	if err != nil {
		return nil, nil, err
	}
	first, last := cidr.AddressRange(ipv4net)
	return first, last, nil
}

func (s *Server) Validate() error {
	return nil
}
