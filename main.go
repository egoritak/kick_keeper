package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//var db *gorm.DB
//var err error
//var client *mongo.Client
var dbURI string

type Kick struct {
	gorm.Model

	CompanyName string  `json:"company"`
	KickName    string  `json:"kickname,omitempty"`
	Status      string  `json:"status,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon, omitempty"`
	Speed       float64 `json:"speed, omitempty"`
}

func main() {
	dial := os.Getenv("DIALECT")
	mydb, _ := getDB(dial)

	/*
		db, err = gorm.Open(dial, dbURI)
		if err != nil {
			log.Fatal("Failed open connection: ", err)
		} else {
			fmt.Println("Successful connection to database")
		}

		//Cose connection
		defer db.Close()

		//Make migrations to the database
		db.AutoMigrate(&Kick{})

				//setup mongodb client
			        ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			        clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
			        client, _ = mongo.Connect(ctx, clientOptions)
	*/

	mydb.SetupConnection()

	//API routes
	//router := mux.NewRouter()
	//router.HandleFunc("/kicks", mydb.PullKicks).Methods("GET")

	//http.ListenAndServe(":8080", router)
}
