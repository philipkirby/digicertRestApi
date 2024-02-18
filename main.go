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
	mongoDSN    = *flag.String("mongoDSN", "mongodb://pdkirby:DontLookAtMyPassword@localhost:27017", "Mongo DSN") // can use program arguments for dsn
	restPort    = *flag.String("restPort", "8081", "rest port")                                                   // can use program arguments for dsn
)

func main() {
	log.Println("Starting rest api")
	if envDSN := os.Getenv("mongoDSN"); envDSN != "" { // override dsn with os environment dsn such as docker dsn
		mongoDSN = envDSN
	}
	if envPort := os.Getenv("restPort"); envPort != "" { // override port with os environment port such as docker dsn
		restPort = envPort
	}

	signal.Notify(closeNotify, os.Kill, os.Interrupt, syscall.SIGTERM) // catch terminate signal to close rest properly

	dbHandler, err := db.CreateMongoDBHandler(mongoDSN)
	if err != nil {
		log.Println(err.Error())
		return
	}

	service, err := internal.CreateRestApiService(dbHandler, restPort)
	if err != nil {
		log.Println(err.Error())
		return
	}

	service.Start()
	<-closeNotify
	log.Println("stopping service")
	service.Stop()
}
