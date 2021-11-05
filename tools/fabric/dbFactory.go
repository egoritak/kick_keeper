package fabric

import (
	"fmt"
	"net/http"

	mongopkg "../../pkg/mongopkg"
	postgrpkg "../../pkg/postgrpkg"
	configReader "../configReader"
)

var config configReader.Config

func readConfig() {
	// for this tutorial, we will hard code it to config.txt
	var err error
	config, err = configReader.ReadConfig(`configs/config.txt`)

	if err != nil {
		fmt.Println(err)
	}
}

var (
	dial     string
	host     string
	port     string
	user     string
	dbname   string
	password string
)

type dbWorker interface {
	PostKick(w http.ResponseWriter, r *http.Request)
	PostKicks(w http.ResponseWriter, r *http.Request)
	Connect()
}

func GetDB(dbType string) (dbWorker, error) {
	//URI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
	//	host, user, dbname, password, port)
	readConfig()

	dial = config["DIALECT"]
	host = config["HOST"]
	port = config["DBPort"]
	user = config["USER"]
	dbname = config["NAME"]
	password = config["PASSWORD"]

	if dbType == "postgres" {
		return newPostgr(), nil
	} else if dbType == "mongodb" {
		return newMongo(), nil
	}

	return nil, fmt.Errorf("Wrong db type passed")
}

func newMongo() dbWorker {
	dbURI := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dial, user, password, host, port, dbname)
	return &mongopkg.MongoDB{
		URI:      dbURI,
		Dialect:  dial,
		Host:     host,
		DBPort:   port,
		User:     user,
		DBName:   dbname,
		Password: password,
	}
}

func newPostgr() dbWorker {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		host, user, dbname, password, port)

	return &postgrpkg.Postgr{
		URI:      dbURI,
		Dialect:  dial,
		Host:     host,
		DBPort:   port,
		User:     user,
		DBName:   dbname,
		Password: password,
	}
}
