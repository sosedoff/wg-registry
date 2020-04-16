package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	strippedConfig, err := exec.Command(c.wgQuickPath, "strip", iface).Output()
	if err != nil {
		log.Println("wg-quick strip failed:", string(strippedConfig))
		return err
	}

	tmpPath := filepath.Join(c.configDir, "temp")
	defer os.Remove(tmpPath)

	if err := ioutil.WriteFile(tmpPath, strippedConfig, 0644); err != nil {
		return err
	}

	out, err := exec.Command(c.wgPath, "syncconf", iface, tmpPath).CombinedOutput()
	if err != nil {
		log.Println("wg-quick strip failed:", string(out))
		return err
	}

	return err
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

	config, err := generate.ServerConfig(c.store, server)
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
