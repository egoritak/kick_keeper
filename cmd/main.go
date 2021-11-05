package main

import (
	configReader "../tools/configReader"
	fabric "../tools/fabric"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/*
type Kick struct {
	//	gorm.Model

	CompanyName string  `json:"company" bson:"company,omitempty"`
	KickName    string  `json:"kickname,omitempty" bson:"kickname,omitempty"`
	Status      string  `json:"status,omitempty" bson:"status,omitempty"`
	Lat         float64 `json:"lat,omitempty" bson:"lat,omitempty"`
	Lon         float64 `json:"lon, omitempty" bson:"lon,omitempty"`
	Speed       float64 `json:"speed, omitempty" bson:"speed,omitempty"`
}
*/
func main() {
	config, _ := configReader.ReadConfig("config/config.txt")
	dial := config["DIALECT"]
	mydb, _ := fabric.GetDB(dial)

	mydb.Connect()
}
