package database

import (
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/spf13/viper"
	_ "github.com/lib/pq"
)

var Db *dbx.DB

func ConnectDbPostgres() {
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	name := viper.GetString("db.name")
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")



	db, err := dbx.MustOpen("postgres", "postgres://"+username+":"+password+"@"+host+":"+port+"/"+name+"?sslmode=disable")

	if (err != nil){
		panic(err)
	}

	Db = db
}