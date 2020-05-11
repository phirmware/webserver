package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when the resource could not be found in the databse
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is returned when an invalid ID is provided to a method like delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

// User defines the shape of our user table
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

// UserService defines the shape of the user service
type UserService struct {
	db *gorm.DB
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// NewUserService returns a userservice type
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

// Close closes the database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// ByID finds a user with provided ID
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail searches user by email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a provided user in the DataBase
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update updates user in the database
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete deletes a user from the databse
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// DestructiveReset drops the user table if it exists and creates another one
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will automatically migrate the user esp for production environment
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
