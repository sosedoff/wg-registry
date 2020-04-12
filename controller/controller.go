package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sosedoff/wg-registry/generate"
	"github.com/sosedoff/wg-registry/store"
)

// Controller handles wireguard network configuration
type Controller struct {
	wgPath      string
	wgQuickPath string
	configDir   string
	store       store.Store
}

// New returns a new controller
func New(wgPath string, wgQuickPath string, configDir string, store store.Store) *Controller {
	return &Controller{
		wgPath:      wgPath,
		wgQuickPath: wgQuickPath,
		configDir:   configDir,
		store:       store,
	}
}

// Restart restarts the network interface
func (c *Controller) Restart(iface string) error {
	if err := c.Down(iface); err != nil {
		return err
	}
	return c.Up(iface)
}

// Reload reloads the interface configuration without a restart
func (c *Controller) Reload(iface string) error {
	configPath := filepath.Join(c.configDir, fmt.Sprintf("%s.conf", iface))

	strippedConfig, err := exec.Command(c.wgQuickPath, "strip", configPath).CombinedOutput()
	if err != nil {
		log.Println("wg-quick strip failed:", string(strippedConfig))
		return err
	}

	syncCmd := exec.Command(c.wgPath, "syncconf", iface)
	syncCmd.Stdin = bytes.NewReader(strippedConfig)

	return syncCmd.Run()
}

// Down brings the network interface down
func (c *Controller) Down(iface string) error {
	out, err := exec.Command(c.wgQuickPath, "down", iface).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "is not a WireGuard interface") {
			err = nil
		}
	}
	return err
}

// Up brings the network interface up
func (c *Controller) Up(iface string) error {
	out, err := exec.Command(c.wgQuickPath, "up", iface).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "already exists") {
			err = nil
		}
	}
	return err
}

// Apply updates server configuration and reconfigures network interface
func (c *Controller) Apply(restart bool) error {
	server, err := c.store.FindServer()
	if err != nil {
		return err
	}
	if server == nil {
		return errors.New("server not found")
	}

	config, err := generate.ServerConfig(c.store)
	if err != nil {
		return err
	}

	configName := fmt.Sprintf("%s.conf", server.Interface)
	configPath := filepath.Join(c.configDir, configName)

	if err := ioutil.WriteFile(configPath, config, 0644); err != nil {
		return err
	}

	if restart {
		return c.Restart(server.Interface)
	}

	return c.Reload(server.Interface)
}
