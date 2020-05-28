package models

import (
	"github.com/jinzhu/gorm"
)

// Services defines the shape of the services
type Services struct {
	Gallery GalleryService
	User    UserService
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
		Gallery: &galleryGorm{},
	}, nil
}
