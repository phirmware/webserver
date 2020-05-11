package main

import (
	"fmt"

	"lenslocked.com/models"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		" dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	// us.DestructiveReset()

	// user := models.User{
	// 	Name:  "chibuzor ojukwu",
	// 	Email: "chibuzorojukwu1@gmail.com",
	// }

	foundUser, err := us.ByEmail("chibuzorojukwu1@gmail.com")
	if err != nil {
		panic(err)
	}
	if err := us.Delete(foundUser.ID); err != nil {
		panic(err)
	}

	_, err = us.ByID(foundUser.ID)
	if err != models.ErrNotFound {
		fmt.Println("Error: User wasnt deleted")
	}

}
