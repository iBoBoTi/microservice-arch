package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iBoBoTi/microservice-arch/authentication-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)


const webPort = "80"

var counts int64

func main(){
	log.Println("Starting Authentication Service...")

	// connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}


	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	// http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	log.Fatal(srv.ListenAndServe())
}

func OpenDB(dsn string) (*sql.DB, error){
	db, err :=  sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database service: %v", err)
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Println(dsn)

	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10{
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}