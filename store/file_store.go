package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sosedoff/wg-registry/model"
)

type FileData struct {
	Users   []model.User   `json:"user"`
	Devices []model.Device `json:"device"`
	Servers []model.Server `json:"servers"`

	UserSeq    int `json:"_users_seq"`
	DeviceSeq  int `json:"_devices_seq"`
	ServersSeq int `json:"_servers_seq"`
}

type FileStore struct {
	sync.Mutex
	path string
	data FileData
}

func NewFileStore(path string) (*FileStore, error) {
	return &FileStore{
		Mutex: sync.Mutex{},
		path:  path,
		data:  FileData{},
	}, nil
}

func (s *FileStore) read() error {
	_, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(s.path, []byte("{}"), 0644); err != nil {
			return err
		}
	}

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.data)
}

func (s *FileStore) write() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.path, data, 0644)
}

func (s *FileStore) mustread() {
	if err := s.read(); err != nil {
		panic(err)
	}
}

func (s *FileStore) mustwrite() {
	if err := s.write(); err != nil {
		panic(err)
	}
}

func (s *FileStore) keyfor(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return strings.ToLower(t.Elem().Name())
	} else {
		return strings.ToLower(t.Name())
	}
}

func (s *FileStore) idval(val interface{}) int {
	switch val.(type) {
	case int:
		return val.(int)
	case string:
		n, err := strconv.Atoi(val.(string))
		if err != nil {
			panic(err)
		}
		return n
	default:
		panic("invalid type")
	}
}

func (s *FileStore) AutoMigrate() error {
	s.Lock()
	defer s.Unlock()

	if err := s.read(); err != nil {
		return err
	}

	if s.data.Users == nil {
		s.data.Users = []model.User{}
	}
	if s.data.Devices == nil {
		s.data.Devices = []model.Device{}
	}
	if s.data.Servers == nil {
		s.data.Servers = []model.Server{}
	}

	return s.write()
}

func (s *FileStore) UserCount() (int, error) {
	return len(s.data.Users), nil
}

func (s *FileStore) FindUserByID(uid interface{}) (*model.User, error) {
	s.Lock()
	defer s.Unlock()

	id := s.idval(uid)

	for _, user := range s.data.Users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (s *FileStore) FindUserByEmail(val string) (*model.User, error) {
	s.Lock()
	defer s.Unlock()

	for _, user := range s.data.Users {
		if user.Email == val {
			return &user, nil
		}
	}
	return nil, nil
}

func (s *FileStore) FindServer() (*model.Server, error) {
	s.Lock()
	defer s.Unlock()

	for _, server := range s.data.Servers {
		return &server, nil
	}

	return nil, nil
}

func (s *FileStore) CreateServer(server *model.Server) error {
	s.Lock()
	defer s.Unlock()

	s.data.ServersSeq++
	server.ID = s.data.ServersSeq

	if len(s.data.Servers) == 0 {
		s.data.Servers = append(s.data.Servers, *server)
	} else {
		s.data.Servers[0] = *server
	}

	return s.write()
}

func (s *FileStore) FindDevice(val interface{}) (*model.Device, error) {
	s.Lock()
	defer s.Unlock()

	id := s.idval(val)

	for _, device := range s.data.Devices {
		if device.ID == id {
			return &device, nil
		}
	}

	return nil, nil
}

func (s *FileStore) FindUserDevice(user *model.User, val interface{}) (*model.Device, error) {
	s.Lock()
	defer s.Unlock()

	id := s.idval(val)

	for _, device := range s.data.Devices {
		if device.UserID == user.ID && device.ID == id {
			return &device, nil
		}
	}

	return nil, nil
}

func (s *FileStore) AllDevices() ([]model.Device, error) {
	s.Lock()
	defer s.Unlock()

	return s.data.Devices, nil
}

func (s *FileStore) AllUsers() ([]model.User, error) {
	s.Lock()
	defer s.Unlock()

	return s.data.Users, nil
}

func (s *FileStore) GetDevicesByUser(val interface{}) ([]model.Device, error) {
	s.Lock()
	defer s.Unlock()

	id := s.idval(val)
	result := []model.Device{}

	for _, device := range s.data.Devices {
		if device.UserID == id {
			result = append(result, device)
		}
	}

	return result, nil
}

func (s *FileStore) CreateUser(user *model.User) error {
	s.Lock()
	defer s.Unlock()

	s.data.UserSeq++

	user.ID = s.data.UserSeq
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	s.data.Users = append(s.data.Users, *user)

	return s.write()
}

func (s *FileStore) CreateDevice(device *model.Device) error {
	if err := device.Validate(); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.data.DeviceSeq++

	device.ID = s.data.DeviceSeq
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()

	s.data.Devices = append(s.data.Devices, *device)

	return s.write()
}

func (s *FileStore) DeleteUserDevice(user *model.User, device *model.Device) error {
	s.Lock()
	defer s.Unlock()

	for idx, d := range s.data.Devices {
		if d.ID == device.ID && d.UserID == device.UserID {
			s.data.Devices = append(s.data.Devices[:idx], s.data.Devices[idx+1:]...)
			break
		}
	}

	return s.write()
}

func (s *FileStore) AllocatedIPV4() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	result := []string{s.data.Servers[0].IPV4Addr}
	for _, d := range s.data.Devices {
		result = append(result, d.IPV4)
	}

	return result, nil
}
