package store

import (
	"errors"
	"net/url"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/sosedoff/wg-registry/model"
)

// Store is an interface for database operations
type Store interface {
	AutoMigrate() error
	UserCount() (int, error)
	FindUserByID(interface{}) (*model.User, error)
	FindUserByEmail(string) (*model.User, error)
	FindServer() (*model.Server, error)
	SaveServer(*model.Server) error
	CreateServer(*model.Server) error
	FindDevice(interface{}) (*model.Device, error)
	FindUserDevice(*model.User, interface{}) (*model.Device, error)
	AllDevices() ([]model.Device, error)
	AllUsers() ([]model.User, error)
	GetDevicesByUser(interface{}) ([]model.Device, error)
	CreateUser(*model.User) error
	CreateDevice(*model.Server, *model.Device) error
	DeleteUserDevice(*model.User, *model.Device) error
	AllocatedIPV4() ([]string, error)
}

// Init returns a new store for the scheme
func Init(connstr string) (Store, error) {
	uri, err := url.Parse(connstr)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "postgresql", "mysql", "sqlite":
		return NewDatabaseStore(uri.Scheme, connstr)
	case "file":
		path := strings.TrimPrefix(connstr, "file://")
		return NewFileStore(path)
	default:
		return nil, errors.New("invalid scheme: " + uri.Scheme)
	}
}

// NextIPV4 returns a new IPV4 allocation
func NextIPV4(store Store, server *model.Server) (string, error) {
	ipfirst, iplast, err := server.IPV4Range()
	if err != nil {
		return "", nil
	}

	allocated, err := store.AllocatedIPV4()
	if err != nil {
		return "", err
	}

	cur := ipfirst
	for {
		taken := false

		for _, val := range allocated {
			if val == cur.String() {
				taken = true
				break
			}
		}

		if !taken {
			break
		}

		if cur.Equal(iplast) {
			return "", errors.New("cant allocate address")
		}

		cur = cidr.Inc(cur)
	}

	return cur.String(), nil
}
