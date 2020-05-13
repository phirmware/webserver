package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
	"lenslocked.com/views"
)

const serverPort = ":8080"

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

var homeView *views.View
var contactView *views.View
var questionsView *views.View

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func home(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Name   string
		Footer string
	}

	pageData := data{
		Name:   "Awesome",
		Footer: "Coutesy",
	}
	setHeaders(w)
	must(contactView.Render(w, pageData))
}

func questions(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	must(questionsView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	fmt.Fprint(w, "<h3>404 Not Found</h3>")
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	// us.DestructiveReset()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)
	galleryC := controllers.NewGallery()

	r := mux.NewRouter()
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.Handle("/FAQ", staticC.FaqView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("GET")
	r.HandleFunc("/login", usersC.HandleLogin).Methods("POST")
	r.HandleFunc("/gallery/new", galleryC.New).Methods("GET")
	r.HandleFunc("/cookie-test", usersC.CookieTest)

	var h http.Handler = http.HandlerFunc(notFound)
	r.NotFoundHandler = h
	fmt.Printf("Server listening on serverPort %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, r))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
