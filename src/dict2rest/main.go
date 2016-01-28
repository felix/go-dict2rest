package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/stretchr/graceful"
	"log"
	"net/http"
	"os"
	"time"
)

// Globals
var dictServer string

func main() {
	var err error

	port := flag.String("port", "8080", "Listen port")
	dictHost := flag.String("dicthost", "localhost", "Dict server name")
	dictPort := flag.String("dictport", "2628", "Dict server port")
	gzip := flag.Bool("gzip", false, "Enable gzip compression")

	flag.Parse()

	dictServer = *dictHost + ":" + *dictPort

	client, err := getDictClient()
	if err != nil {
		os.Exit(1)
	}

	// Get the global list of databases
	dicts, err := getDictionaries(client)
	if err != nil {
		log.Println("Unable to get database list")
		os.Exit(1)
	}

	for _, d := range dicts {
		log.Printf("Available dictionary %s: %s", d.Name, d.Desc)
	}
	// No need for it until a request comes in
	client.Close()

	// Define our routes
	router := httprouter.New()
	router.GET("/define/:word", dictDefine)
	router.GET("/databases", dictDatabases)
	router.GET("/db", dictDatabases)

	// Define our middlewares

	// Going to need some CORS headers
	cors := cors.New(cors.Options{
		AllowedHeaders: []string{
			"Accept", "Content-Type", "Origin",
		},
		AllowedMethods: []string{
			"GET", "OPTIONS",
		},
	})

	chain := alice.New(cors.Handler, Logger).Then(router)
	if *gzip {
		chain = alice.New(cors.Handler, Logger, Gzip).Then(router)
		log.Println("Using Gzip compression")
	}

	listen := ":" + *port

	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server:  &http.Server{Addr: listen, Handler: chain},
	}

	log.Printf("Listening at %s", listen)
	srv.ListenAndServe()
}
