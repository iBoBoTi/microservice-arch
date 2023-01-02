package main

import (
	"fmt"
	"log"
	"net/http"
)


const webPort = "80"

func main(){
	log.Println("Starting Authentication Service...")

	//TODO: connect to db
	

	app := Config{}

	// http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	log.Fatal(srv.ListenAndServe())
}