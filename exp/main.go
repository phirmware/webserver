package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	g, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}
	err = g.Create(&user).Error
	if err != nil {
		panic(err)
	}
	// Verify that the user has a Remember and RememberHash
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}
	// Now verify that we can lookup a user with that remember
	// token
	// user2, err := us.ByRemember(user.Remember)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", *user2)

}
