package models

import (
	"github.com/jinzhu/gorm"
)

// Gallery defines the shape of the gallery table in our db
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// GalleryService defines the shape of gallery service interface
type GalleryService interface {
	GalleryDB
}

// GalleryDB defines the shape of the gallery db interface
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return nil
}
