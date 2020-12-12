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
)

const serverPort = ":8080"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "lenslocked_dev"
	password = "password"
)

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	fmt.Fprint(w, "<h3>404 Not Found</h3>")
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	if err := services.AutoMigrate(); err != nil {
		panic(err)
	}
	// services.DestructiveReset()

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	userMw := middleware.User{
		UserService: services.User,
	}

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.Handle("/FAQ", staticC.FaqView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("GET")
	r.HandleFunc("/login", usersC.HandleLogin).Methods("POST")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/cookie-test", usersC.CookieTest)

	var h http.Handler = http.HandlerFunc(notFound)
	r.NotFoundHandler = h
	fmt.Printf("Server listening on serverPort %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, userMw.Apply(r)))
}
