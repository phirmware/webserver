package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"

	"lenslocked.com/views"
)

// Users defines the users controller
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// SignupForm defines the shape of the users input
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm defines the shape of the login form
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUsers Renders our new users page
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New is the controller for GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Login handles the GET /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	if err := u.LoginView.Render(w, nil); err != nil {
		panic(err)
	}
}

// .x3RFTVjGyA9SMysDJVtDW6ttZK

// Create is the controller that creates a new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookie-test", http.StatusFound)
}

// HandleLogin handles the request to login from the user
func (u *Users) HandleLogin(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrPasswordIncorrect:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookie-test", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest as the name suggests
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusOK)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusOK)
		return
	}
	fmt.Fprintln(w, user)
}
