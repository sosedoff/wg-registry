package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sosedoff/x/wireguard-manager/generate"
	"github.com/sosedoff/x/wireguard-manager/store"
)

// Controller handles wireguard network configuration
type Controller struct {
	wgQuickPath string
	configDir   string
	store       store.Store
	updates     chan struct{}
	stop        chan struct{}
}

// New returns a new controller
func New(wgQuickPath string, configDir string, store store.Store) *Controller {
	return &Controller{
		wgQuickPath: wgQuickPath,
		configDir:   configDir,
		store:       store,
		updates:     make(chan struct{}),
		stop:        make(chan struct{}),
	}
}

// ScheduleApply queues the configuration change
func (c *Controller) ScheduleApply() {
	c.updates <- struct{}{}
}

// Start start the controller worker
func (c *Controller) Start() {
	log.Println("starting controller")
	defer log.Println("stopped controller")

	for {
		select {
		case <-c.updates:
			log.Println("applying configuration")
			if err := c.Apply(); err != nil {
				log.Println("config apply error:", err)
			}

		case <-c.stop:
			log.Println("stopping controller")
			return
		}
	}
}

// Stop stops the controller worker
func (c *Controller) Stop() {
	close(c.updates)
	c.stop <- struct{}{}
}

// Restart restarts the network interface
func (c *Controller) Restart(iface string) error {
	if err := c.Down(iface); err != nil {
		return err
	}
	return c.Up(iface)
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
func (c *Controller) Apply() error {
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

	return c.Restart(server.Interface)
}
