package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"

	"lenslocked.com/views"
)

// Users defines the users controller
type Users struct {
	NewView *views.View
	us      *models.UserService
}

// SignupForm defines the shape of the users input
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUsers Renders our new users page
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

// New is the controller for GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is the controller that creates a new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
	fmt.Fprintln(w, "Name is", form.Name)
}
