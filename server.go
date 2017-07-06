package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"database/sql"

	"github.com/sfreiberg/gotwilio"
)

// Configuration : Core config structure
type Configuration struct {
	Twilio TwilioKey
	DB DbCreds
}

// TwilioKey : Config strucuture for Twilio
type TwilioKey struct {
	Sid    string
	Token  string
	Number string
}

type DbCreds struct {
	Host string
	Port int
	User string
	Pass string
	DbName string
}

// Globals
var (
	done           = make(chan struct{})
	configuration  = loadConfig()
	antidoseTwilio = loadTwilio()
	db             = loadDB()
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
}

func failGracefully(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

func loadConfig() Configuration {
	file, err := os.Open("./config/conf.json")
	failOnError(err, "Config json not found. Make sure it is present.")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func loadTwilio() *gotwilio.Twilio {
	return gotwilio.NewTwilioClient(configuration.Twilio.Sid, configuration.Twilio.Token)
}

func loadDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configuration.DB.Host, configuration.DB.Port, configuration.DB.User, configuration.DB.Pass, configuration.DB.DbName)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	failOnError(err, "Failed to open Postgres")

	err = db.Ping()
	failOnError(err, "Failed to ping Postgres")
	return db
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	initRoutes()
	<-done
}
