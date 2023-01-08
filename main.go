package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Kavuti/goauth/roles"
	"github.com/Kavuti/goauth/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Database Connection
	log.Println("Connecting to database")
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	check(err)
	defer db.Close()

	// Database Migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	check(err)
	log.Println("Applying migrations to database")
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	check(err)
	err = m.Up()
	if err != nil && reflect.TypeOf(err).String() == "migrate.ErrDirty" {
		fmt.Printf("%#v\n", err)
		if os.Getenv("FORCE_MIGRATION") == "true" {
			m.Force(err.(migrate.ErrDirty).Version)
		}
	}
	if err != nil && err != migrate.ErrNoChange {
		check(err)
	}
	log.Println("Database migrations applied. Starting the service")

	// Router Configuration
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Handlers registration
	usersHandler := users.NewUserHandler(r, db)
	rolesHandler := roles.NewRoleHandler(r, db)

	r.Mount("/users", usersHandler.Routes())
	r.Mount("/roles", rolesHandler.Routes())

	// Server start listening
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server started on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
