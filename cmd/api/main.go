package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// set application config
	var app application

	// Loading env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// read from command line
	app.DSN = os.Getenv("APP_DSN")
	app.DSN = strings.Replace(app.DSN, "localhost", "postgres", 1)

	app.JWTSecret = os.Getenv("APP_JWTSECRET")
	app.JWTIssuer = os.Getenv("APP_JWTISSUER")
	app.JWTAudience = os.Getenv("APP_JWTAUDIENCE")
	app.CookieDomain = os.Getenv("APP_COOKIE_DOMAIN")
	app.Domain = os.Getenv("APP_DOMAIN")
	app.APIKey = os.Getenv("APP_API_KEY")

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "_Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
