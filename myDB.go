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
	router.HandleFunc("/post/kick", mydb.PostKick).Methods("POST")
	router.HandleFunc("/post/kicks", mydb.PostKicks).Methods("POST")

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
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		clientOptions := options.Client().ApplyURI("mongodb://postgres:qwerty@localhost:27017/kick_keeper") //"mongodb://localhost:27017")
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal("Failed connect to mongodb", err)
		} else {
			fmt.Println("Successful connection to database")
		}
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
		collection := client.Database(mydb.DBName).Collection("kicks")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(w).Encode(kicks)
	}
}

func (mydb myDB) PostKick(w http.ResponseWriter, r *http.Request) {
	if mydb.Dialect == "postgres" {
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

	} else if mydb.Dialect == "mongodb" {
		w.Header().Set("content-type", "application/json")
		var kick Kick
		err := json.NewDecoder(r.Body).Decode(&kick)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Bad Json")
			return
		}
		collection := client.Database(mydb.DBName).Collection("kicks")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, err := collection.InsertOne(ctx, kick)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	}
}
func (mydb myDB) PostKicks(w http.ResponseWriter, r *http.Request) {
	if mydb.Dialect == "postgres" {
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
	} else if mydb.Dialect == "mongodb" {
		w.Header().Set("content-type", "application/json")
		var kicks []Kick
		err := json.NewDecoder(r.Body).Decode(&kicks)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Bad Json")
			return
		}
		collection := client.Database(mydb.DBName).Collection("kicks")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		for idx := range kicks {
			result, err := collection.InsertOne(ctx, kicks[idx])
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
			}
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				w.WriteHeader(http.StatusConflict)
			}
		}
		json.NewEncoder(w).Encode("All kicks were pushed")
	}
}
