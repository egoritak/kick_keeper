package main

import (
	"fmt"
	"os"
)

type postgr struct {
	myDB
}

func newPostgr() dbWorker {
	dial := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbURI = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		host, user, dbname, password, port)
	return &postgr{
		myDB: myDB{
			Dialect:  dial,
			Host:     host,
			DBPort:   port,
			User:     user,
			DBName:   dbname,
			Password: password,
		},
	}
}
