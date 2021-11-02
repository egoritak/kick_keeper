package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbURI string

type Kick struct {
	gorm.Model

	CompanyName string  `json:"company" bson:"company,omitempty"`
	KickName    string  `json:"kickname,omitempty" bson:"kickname,omitempty"`
	Status      string  `json:"status,omitempty" bson:"status,omitempty"`
	Lat         float64 `json:"lat,omitempty" bson:"lat,omitempty"`
	Lon         float64 `json:"lon, omitempty" bson:"lon,omitempty"`
	Speed       float64 `json:"speed, omitempty" bson:"speed,omitempty"`
}

func main() {
	dial := os.Getenv("DIALECT")
	mydb, _ := getDB(dial)

	mydb.SetupConnection()
}
