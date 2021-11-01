package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type myDB struct {
	Dialect  string
	Host     string
	DBPort   string
	User     string
	DBName   string
	Password string
}

var db *gorm.DB
var err error
var client *mongo.Client

func (mydb myDB) SetupConnection() error {

	//API routes
	router := mux.NewRouter()
	router.HandleFunc("/kicks", mydb.PullKicks).Methods("GET")

	if mydb.Dialect == "postgres" {
		db, err = gorm.Open(mydb.Dialect, dbURI)
		if err != nil {
			log.Fatal("Failed open connection: ", err)
		} else {
			fmt.Println("Successful connection to database")
		}

		//Make migrations to the database
		db.AutoMigrate(&Kick{})

		ConnectPsql(router)

	} else if mydb.Dialect == "mongodb" {
		//setup mongodb client
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, _ = mongo.Connect(ctx, clientOptions)
		ConnectMongo(router)
	}

	return nil
}

func ConnectPsql(router http.Handler) {
	//close connection
	defer db.Close()
	http.ListenAndServe(":8080", router)
}

func ConnectMongo(router http.Handler) {
	http.ListenAndServe(":8080", router)
}

func (mydb myDB) PullKicks(w http.ResponseWriter, r *http.Request) {
	if mydb.Dialect == "postgres" {
		var kicks []Kick
		db.Find(&kicks)

		json.NewEncoder(w).Encode(&kicks)

	} else if mydb.Dialect == "mongodb" {
		w.Header().Set("content-type", "application/json")
		var kicks []Kick
		collection := client.Database("mongodb").Collection("kicks")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var kick Kick
			cursor.Decode(&kick)
			kicks = append(kicks, kick)
		}
		if err := cursor.Err(); err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(w).Encode(kicks)
	}
}

func (mydb myDB) PullKick(w http.ResponseWriter, r *http.Request) {
}
func (mydb myDB) PostKick(w http.ResponseWriter, r *http.Request) {
}
func (mydb myDB) PostKicks(w http.ResponseWriter, r *http.Request) {
}
func (mydb myDB) DeleteKick(w http.ResponseWriter, r *http.Request) {
}
func (mydb myDB) DeleteKicks(w http.ResponseWriter, r *http.Request) {
}
