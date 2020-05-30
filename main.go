package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/middleware"
	"lenslocked.com/models"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func notFound(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	fmt.Fprint(w, "<h3>404 Not Found</h3>")
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset()

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	r := mux.NewRouter()
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.Handle("/FAQ", staticC.FaqView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("GET")
	r.HandleFunc("/login", usersC.HandleLogin).Methods("POST")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/cookie-test", usersC.CookieTest)

	var h http.Handler = http.HandlerFunc(notFound)
	r.NotFoundHandler = h
	fmt.Printf("Server listening on serverPort %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, r))
}
