package store

import (
	"errors"
	"net/url"
	"strings"

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
	CreateDevice(*model.Device) error
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
