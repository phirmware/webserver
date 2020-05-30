package models

import (
	"github.com/jinzhu/gorm"
)

// Services defines the shape of the services
type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// NewServices returns our services struct
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		db:      db,
	}, nil
}

// Close closes all connection to the database
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate creates the tables in our db
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// DestructiveReset drops the tables in the database and rebuilds
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
