#! /usr/bin/env bash

docker rmi webdevgo

docker build -f dockerfile -t webdevgo .

docker run --rm -d  -p 5432:5432 -e POSTGRES_DB=lenslocked_dev webdevgo:latest

echo "Run GO file main.go TO CHECK IF DB WAS SUCCESSFULLY CREATED"