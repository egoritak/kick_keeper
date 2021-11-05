# Kick Keeper

Kick Keeper is a service that allows you to store information about scooters in the Mongo and PostgreSQL databases.

## Setup
Write configuration in .env file to use your database.
.env file contains:
- DIALECT 
Dialect is the database you use. It can be "mongodb" or "postgres".
- HOST
Host contains address of database. It's "localhost" by default.
- DBPORT
DBPort - is the port using to connect to database. It's "5432" by default.
- USER
In this field you have to write the "username" of your database.
- NAME
This field contains the name of oppened database.
- PASSWORD
This field contains password to get access to database.


## Usage
After you have finished setup, write the following line in console to be able to read this settings from the programm:
```sh
source .env
```
Then launch service by writing:
```sh
go run .
```

Now you can send POST json requests.

## HANDLERS
"/kick" - to post the only json document.
"/kicks" - to post many json documents.

## JSON structure
{<br />
    "company": string,<br />
    "kickname": string,<br />
    "status": string,<br />
    "lat": float64,<br />
    "lon": float64, <br />
    "speed": float64 <br />
} 

## Packages
Run:
```sh
go get "github.com/jinzhu/gorm"
go get "github.com/jinzhu/gorm/dialects/postgres"
go get "github.com/gorilla/mux"
go get "github.com/jinzhu/gorm"
go get "go.mongodb.org/mongo-driver/bson"
go get "go.mongodb.org/mongo-driver/mongo"
go get "go.mongodb.org/mongo-driver/mongo/options"
```
