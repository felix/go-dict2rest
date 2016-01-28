package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/stretchr/graceful"
	"golang.org/x/net/dict"
	"log"
	"net/http"
	"os"
	"time"
)

var client *dict.Client
var dictList map[string]dict.Dict

func main() {
	var err error

	port := flag.String("port", "8080", "Listen port")
	dictHost := flag.String("dicthost", "localhost", "Dict server name")
	dictPort := flag.String("dictport", "2628", "Dict server port")
	gzip := flag.Bool("gzip", false, "Enable gzip compression")

	flag.Parse()

	server := *dictHost + ":" + *dictPort
	client, err = dict.Dial("tcp", server)
	if err != nil {
		log.Printf("Unable to connect to dict server %s", server)
		os.Exit(1)
	}
	log.Println("Connected to", server)

	var dictArr []dict.Dict
	dictArr, err = client.Dicts()
	if err != nil {
		log.Fatal("Unable to get dictionaries")
	}

	dictList = make(map[string]dict.Dict)
	for _, d := range dictArr {
		log.Println("Using dictionary", d.Name, d.Desc)
		dictList[d.Name] = d
	}

	router := httprouter.New()
	//router.GET("/", Index)
	router.GET("/define/:word", Define)

	chain := alice.New(Logger).Then(router)
	if *gzip {
		chain = alice.New(Logger, Gzip).Then(router)
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
