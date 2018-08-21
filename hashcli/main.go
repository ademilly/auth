package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ademilly/auth"
)

var (
	username string
	password string
)

func hash(username, password string) (auth.User, error) {
	newUser, err := auth.Hash(auth.User{Username: username, Password: password})
	if err != nil {
		return auth.User{}, fmt.Errorf("could not encrypt password: %v", err)
	}

	return newUser, nil
}

func main() {
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.Parse()

	if username == "" || password == "" {
		log.Fatalln("username and/or password flag is not set")
	}

	user, err := hash(username, password)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user.Hash)
}
