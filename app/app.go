package app

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Quantaly/mltea-reviews/app/db"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	log       *log.Logger
	templates *template.Template
	db        *pgxpool.Pool
}

func New(log *log.Logger, databaseURL string) (*App, error) {
	a := new(App)
	a.log = log
	err := a.init(databaseURL)
	if err != nil {
		// a.log is def set up
		a.log.Println(err)
		return nil, err
	} else {
		return a, nil
	}
}

func (a *App) init(databaseURL string) (err error) {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return errors.New("DATABASE_URL environment variable not set")
	}

	a.templates = template.New("blank")
	_, err = a.templates.ParseGlob("web/templates/*/*.html")
	if err != nil {
		return
	}
	a.log.Println(a.templates.DefinedTemplates())

	a.db, err = db.SetupConnectionPool(context.Background(), databaseURL)
	if err != nil {
		return
	}

	return nil
}

func (a *App) Run(listenAddr string) error {
	r := mux.NewRouter()
	r.Use(handlers.CompressHandler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.HandleFunc("/", a.getIndex).Methods("GET")
	r.HandleFunc("/teas", a.getTeas).Methods("GET")
	r.HandleFunc("/reviews", a.getReviews).Methods("GET")

	r.HandleFunc("/review", a.postReview).Methods("POST")

	a.log.Println("Listening on", listenAddr)
	err := http.ListenAndServe(listenAddr, r)
	a.log.Println(err)
	a.db.Close()
	return err
}
