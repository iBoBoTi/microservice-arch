package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main(){
	app := Config{}

	log.Printf("Starting broker service at port: %s\n", webPort)

	// http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	log.Fatal(srv.ListenAndServe())
}

