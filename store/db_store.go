package store

import (
	"github.com/jinzhu/gorm"

	"github.com/sosedoff/x/wireguard-manager/model"
)

type DbStore struct {
	db *gorm.DB
}

func NewDatabaseStore(scheme, connstr string) (*DbStore, error) {
	db, err := gorm.Open(scheme, connstr)
	if err != nil {
		return nil, err
	}
	return &DbStore{db: db}, nil
}

func (s *DbStore) AutoMigrate() error {
	return s.db.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.Server{},
	).Error
}

func (s *DbStore) UserCount() (int, error) {
	count := 0
	err := s.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (s *DbStore) FindUserByID(id interface{}) (*model.User, error) {
	user := &model.User{}
	err := s.db.First(user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		user = nil
		err = nil
	}
	return user, err
}

func (s *DbStore) FindUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := s.db.First(user, "LOWER(email) = ?", email).Error
	if err == gorm.ErrRecordNotFound {
		user = nil
		err = nil
	}
	return user, err
}

func (s *DbStore) FindServer() (*model.Server, error) {
	server := &model.Server{}
	err := s.db.First(server).Error
	if err == gorm.ErrRecordNotFound {
		server = nil
		err = nil
	}
	return server, err
}

func (s *DbStore) CreateServer(server *model.Server) error {
	if err := server.Validate(); err != nil {
		return err
	}
	return s.db.Create(server).Error
}

func (s *DbStore) FindDevice(id interface{}) (*model.Device, error) {
	result := &model.Device{}
	err := s.db.Where("id = ?", id).Find(&result).Error
	return result, err
}

func (s *DbStore) FindUserDevice(user *model.User, id interface{}) (*model.Device, error) {
	result := &model.Device{}
	err := s.db.
		Where("user_id = ? AND id = ?", user.ID, id).
		First(&result).
		Error
	if err == gorm.ErrRecordNotFound {
		result = nil
	}
	return result, err
}

func (s *DbStore) AllDevices() ([]model.Device, error) {
	result := []model.Device{}
	err := s.db.Find(&result).Error
	return result, err
}

func (s *DbStore) AllUsers() ([]model.User, error) {
	result := []model.User{}
	err := s.db.Find(&result).Error
	return result, err
}

func (s *DbStore) GetDevicesByUser(id interface{}) ([]model.Device, error) {
	result := []model.Device{}
	err := s.db.Where("user_id = ?", id).Find(&result).Error
	return result, err
}

func (s *DbStore) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *DbStore) CreateDevice(device *model.Device) error {
	if err := device.Validate(); err != nil {
		return err
	}
	return s.db.Create(device).Error
}

func (s *DbStore) DeleteUserDevice(user *model.User, device *model.Device) error {
	return s.db.Delete(device).Error
}

func (s *DbStore) AllocatedIPV4() ([]string, error) {
	addrs := []string{}
	err := s.db.Model(&model.Device{}).Pluck("ip_v4", &addrs).Error
	return addrs, err
}
