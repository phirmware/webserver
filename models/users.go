package models

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"

	"github.com/jinzhu/gorm"
)

const (
	// ErrNotFound is returned when the resource could not be found in the databse
	ErrNotFound modelError = "models: resource not found"
	// ErrIDInvalid is returned when an invalid ID is provided to a method like delete
	ErrIDInvalid modelError = "models: ID provided was invalid"
	userPwPepper            = "secret-random-string"
	// ErrPasswordIncorrect is returned for invalid authentication from a user
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrPasswordTooShort is returned when a user inputs a password that is less than 8 characters
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrPasswordRequired is returned when a password is not provided
	ErrPasswordRequired modelError = "models: password is required"
	// ErrEmailRequired is returned when an email address isnt provided
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when the provided email does not match our requirements from regex
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is retured when the provided email address is already in use
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrRememberRequired is returned when a create or update is attempted without a user remember token hash
	ErrRememberRequired modelError = "models: remember token is required"
	// ErrRememberTooShort is required when a remember token is no at least 32 bytes
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
)

const hmacSecretKey = "secret-hmac-key"

// User defines the shape of our user table
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userService defines the shape of the user service
type userService struct {
	UserDB
}

// UserDB is used to interact with the database
type UserDB interface {
	// Methods for quering for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// userGorm represents our database implementation layer
type userGorm struct {
	db *gorm.DB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

type userValFn func(*User) error

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

var _ UserDB = &userGorm{}
var _ UserService = &userService{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user and back fill data
func (uv *userValidator) Create(user *User) error {
	if err := runUserValFns(user, uv.passwordRequired, uv.passwordMinLength, uv.bcryptPassword, uv.passwordHashRequired, uv.setRememberIfUnset, uv.rememberMinBytes, uv.hmacRemember, uv.rememberHashRequired, uv.normalizeEmail, uv.requireEmail, uv.emailFormat, uv.emailIsAvail); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// ByEmail will normalize an email address before passing it on to the database layer to perform the query
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// NewUserService returns a user service type
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
}

// ByID finds a user with provided ID
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail searches user by email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a provided user in the DataBase
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID < n {
			return ErrIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	// return the error if we got one apart from errNotFound
	if err != nil {
		return err
	}

	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

// A util function that loops over validation functions and executes in the input order
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// Authenticate logins the user
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}

// Update updates user in the database
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Update runs validations and calls the Update method on the next UserDB type
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFns(user, uv.passwordMinLength, uv.bcryptPassword, uv.passwordHashRequired, uv.rememberMinBytes, uv.hmacRemember, uv.rememberHashRequired, uv.normalizeEmail, uv.requireEmail, uv.emailFormat, uv.emailIsAvail); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// ByRemember looks up a user by a given remember token and returns the user and an error(nil/value)
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByRemember will hash the remmeber token and then call the ByRemember on the next UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Delete deletes a user from the databse
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Delete runs validations and passes to the next Delete on UserDB interface
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}
