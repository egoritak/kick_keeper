package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Kick struct {
	gorm.Model

	CompanyName string  `json:"company"`
	KickName    string  `json:"kickname,omitempty"`
	Status      string  `json:"status,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon, omitempty"`
	Speed       float64 `json:"speed, omitempty"`
}

var db *gorm.DB
var err error

type DBWorker struct {
	Dialect  string
	Host     string
	DBPort   string
	User     string
	DBName   string
	Password string
}

//API controllers
func (dbWorker *DBWorker) PullKicks(w http.ResponseWriter, r *http.Request) {
	var kicks []Kick
	db.Find(&kicks)

	json.NewEncoder(w).Encode(&kicks)
}

func (dbWorker *DBWorker) PullKick(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var kick Kick

	db.First(&kick, params["id"])

	json.NewEncoder(w).Encode(kick)
}

func (dbWorker *DBWorker) PostKick(w http.ResponseWriter, r *http.Request) {
	var kick Kick

	err = json.NewDecoder(r.Body).Decode(&kick)

	if err != nil {
		json.NewEncoder(w).Encode("Json Error")
		return
	}

	createdKick := db.Create(&kick)
	err = createdKick.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&kick)
	}
}

func (dbWorker *DBWorker) PostKicks(w http.ResponseWriter, r *http.Request) {
	var kicks []Kick

	err = json.NewDecoder(r.Body).Decode(&kicks)

	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	for idx := range kicks {
		createdKick := db.Create(&kicks[idx])
		err = createdKick.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}
	}
	json.NewEncoder(w).Encode("All kicks were pushed")
}

func (dbWorker *DBWorker) DeleteKick(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var kick Kick

	db.First(&kick, params["id"])
	db.Delete(&kick)

	json.NewEncoder(w).Encode(&kick)
}

func (dbWorker *DBWorker) DeleteKicks(w http.ResponseWriter, r *http.Request) {
	var kicks []Kick
	db.Delete(&kicks)
	json.NewEncoder(w).Encode("All kicks was deleted")
}

func main() {

	dbWorker := DBWorker{}
	//loading enviroment variables
	dbWorker.Dialect = os.Getenv("DIALECT")
	dbWorker.Host = os.Getenv("HOST")
	dbWorker.DBPort = os.Getenv("DBPORT")
	dbWorker.User = os.Getenv("USER")
	dbWorker.DBName = os.Getenv("NAME")
	dbWorker.Password = os.Getenv("PASSWORD")

	//Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		dbWorker.Host, dbWorker.User, dbWorker.DBName, dbWorker.Password, dbWorker.DBPort)

	db, err = gorm.Open(dbWorker.Dialect, dbURI)
	if err != nil {
		log.Fatal("Failed open connection: ", err)
	} else {
		fmt.Println("Successful connection to database")
	}

	//Cose connection
	defer db.Close()

	//Make migrations to the database
	db.AutoMigrate(&Kick{})

	//API routes
	router := mux.NewRouter()

	router.HandleFunc("/kicks", dbWorker.PullKicks).Methods("GET")
	router.HandleFunc("/kick/{id}", dbWorker.PullKick).Methods("GET")

	router.HandleFunc("/post/kick", dbWorker.PostKick).Methods("POST")
	router.HandleFunc("/post/kicks", dbWorker.PostKicks).Methods("POST")

	router.HandleFunc("/delete/kick/{id}", dbWorker.DeleteKick).Methods("DELETE")
	router.HandleFunc("/delete/kicks", dbWorker.DeleteKicks).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
