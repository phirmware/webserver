package models

import (
	"github.com/jinzhu/gorm"
)

const (
	// ErrUserIDRequired is returned when a userID isn't provided
	ErrUserIDRequired modelError = "models: user ID is required"
	// ErrTitleRequired is returned when a title isn't provided
	ErrTitleRequired modelError = "models: title is required"
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
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

type galleryValidator struct {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type galleryValFn func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

// NewGalleryService returns the galleryService struct
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	if err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	); err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	if err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	); err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}
