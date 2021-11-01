package main

import "fmt"

func getDB(dbType string) (dbWorker, error) {
	if dbType == "postgres" {
		return newPostgr(), nil
	} else if dbType == "mongo" {
		return newMongo(), nil
	}

	return nil, fmt.Errorf("Wrong db type passed")
}