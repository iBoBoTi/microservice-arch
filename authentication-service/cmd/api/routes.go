package main

import (
	"net/http"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iBoBoTi/microservice-arch/authentication-service/data"
)

type Config struct {
	DB *sql.DB
	Models data.Models
}

func (app *Config) routes() http.Handler{
	mux := chi.NewRouter()

	// who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposedHeaders: []string{"link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", app.Authenticate)
	
	return mux
} 