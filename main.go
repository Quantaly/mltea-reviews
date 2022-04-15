package main

import (
	"log"
	"os"

	"github.com/Quantaly/mltea-reviews/app"
)

func main() {
	_, heroku := os.LookupEnv("DYNO")

	// set up logger -- if running on heroku, don't need to include timestamp
	logger := log.New(os.Stderr, "", log.LstdFlags)
	if heroku {
		logger.SetFlags(log.Lshortfile)
	} else {
		logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Fatalln("PORT environment variable not set")
	}

	// if running on heroku, bind on all interfaces, else bind only to loopback
	var listenAddr string
	if heroku {
		listenAddr = ":" + port
	} else {
		listenAddr = "127.0.0.1:" + port
	}

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		logger.Fatalln("DATABASE_URL environment variable not set")
	}

	a, err := app.New(logger, databaseURL)
	if err == nil {
		a.Run(listenAddr)
	}
	os.Exit(1)
}
