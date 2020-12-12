package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "lenslocked_dev"
	password = "password"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		" dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	_, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Postgres Database has been setup successfully")
}
