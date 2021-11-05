package postgrpkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kickpkg "../kickpkg"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Postgr struct {
	URI      string
	Dialect  string
	Host     string
	DBPort   string
	User     string
	DBName   string
	Password string
}

var db *gorm.DB

func (myPostgr Postgr) Connect() {
	//API routes
	router := mux.NewRouter()
	router.HandleFunc("/kick", myPostgr.PostKick).Methods("POST")
	router.HandleFunc("/kicks", myPostgr.PostKicks).Methods("POST")

	var err error
	db, err = gorm.Open(myPostgr.Dialect, myPostgr.URI)
	if err != nil {
		log.Fatal("Failed open connection: ", err)
	} else {
		fmt.Println("Successful connection to database")
	}

	//Make migrations to the database
	db.AutoMigrate(&kickpkg.Kick{})

	defer db.Close()

	http.ListenAndServe(":8080", router)
}

func (myPostgr Postgr) PostKick(w http.ResponseWriter, r *http.Request) {
	var kick kickpkg.Kick

	err := json.NewDecoder(r.Body).Decode(&kick)

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

func (myPostgr Postgr) PostKicks(w http.ResponseWriter, r *http.Request) {
	var kicks []kickpkg.Kick

	err := json.NewDecoder(r.Body).Decode(&kicks)

	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	db, err = gorm.Open(myPostgr.Dialect, myPostgr.URI)
	for idx := range kicks {
		createdKick := db.Create(&kicks[idx])
		err = createdKick.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}
	}
	json.NewEncoder(w).Encode("All kicks were pushed")
}
