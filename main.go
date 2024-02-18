package main

import (
	"dockerrestapi/db"
	"dockerrestapi/internal"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	closeNotify = make(chan os.Signal)
	mongoDSN    = *flag.String("mongoDSN", "mongodb://pdkirby:DontLookAtMyPassword@localhost:27017", "Mongo DSN") // can use program arugments for dsn
)

func main() {
	signal.Notify(closeNotify, os.Kill, os.Interrupt, syscall.SIGTERM)

	if envDSN := os.Getenv("mongoDSN"); envDSN != "" { // overide dsn with os enviroment dsn such as docker dsn
		mongoDSN = envDSN
	}
	log.Println("starting rest api")
	conn, err := db.CreateMongoDBHandler(mongoDSN)
	if err != nil {
		log.Println(err.Error())
		return
	}

	service, err := internal.CreateRestApiService(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	service.Start()
	<-closeNotify
	log.Println("stopping service")
	service.Stop()
}
