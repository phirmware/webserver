#! /usr/bin/env bash

function checkDB() {
    go run main.go
    if [ $? -eq 0 ]; then
       echo "DB was successfully setup, Happy coding"
    else
       echo "Something went wrong, try running go file"
    fi
}

docker rmi webdevgo
if [ $? -eq 0 ]; then
   echo "Successfully removed Image"
else
   echo "Image not found, creating image"
fi

docker build -f dockerfile -t webdevgo .
if [ $? -eq 0 ]; then
   echo "Successfully Built docker image"
else
   echo "Something went wrong"
fi

docker run --rm -d  -p 5432:5432 -e POSTGRES_DB=lenslocked_dev -e POSTGRES_PASSWORD=password webdevgo:latest
if [ $? -eq 0 ]; then
   echo "Postgres Container running succesfully"
   echo "Checking db connectivity, please wait........."
   sleep 4
   checkDB
else
   echo "Something went wrong"
fi
