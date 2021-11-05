package mongopkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	kickpkg "../kickpkg"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	URI      string
	Dialect  string
	Host     string
	DBPort   string
	User     string
	DBName   string
	Password string
}

var client *mongo.Client

func (myMongo MongoDB) Connect() {
	//setup mongodb client
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	fmt.Println(myMongo.URI)
	clientOptions := options.Client().ApplyURI(myMongo.URI)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed connect to mongodb", err)
	} else {
		fmt.Println("Successful connection to database")
	}

	//API routes
	router := mux.NewRouter()
	router.HandleFunc("/kick", myMongo.PostKick).Methods("POST")
	router.HandleFunc("/kicks", myMongo.PostKicks).Methods("POST")
	router.HandleFunc("/pull/kicks", myMongo.PullKicks).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func (myMongo MongoDB) PostKick(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var kick kickpkg.Kick

	err := json.NewDecoder(r.Body).Decode(&kick)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad Json")
		return
	}

	collection := client.Database(myMongo.DBName).Collection("kicks")
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

func (myMongo MongoDB) PostKicks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var kicks []kickpkg.Kick
	err := json.NewDecoder(r.Body).Decode(&kicks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad Json")
		return
	}

	collection := client.Database(myMongo.DBName).Collection("kicks")
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

func (myMongo MongoDB) PullKicks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var kicks []kickpkg.Kick
	collection := client.Database(myMongo.DBName).Collection("kicks")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var kick kickpkg.Kick
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
