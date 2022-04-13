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
	"github.com/jackc/pgx/v4"
)

type App struct {
	log        *log.Logger
	isHeroku   bool
	listenAddr string
	templates  *template.Template
	db         *pgx.Conn
}

func New() (*App, error) {
	a := new(App)
	err := a.init()
	if err != nil {
		// a.log is def set up
		a.log.Println(err)
		return nil, err
	} else {
		return a, nil
	}
}

func (a *App) init() (err error) {
	a.log = log.New(os.Stderr, "", log.LstdFlags)

	_, a.isHeroku = os.LookupEnv("DYNO")
	if a.isHeroku {
		a.log.SetFlags(log.Lshortfile)
	} else {
		a.log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		return errors.New("PORT environment variable not set")
	}
	if a.isHeroku {
		a.listenAddr = ":" + port
	} else {
		a.listenAddr = "127.0.0.1:" + port
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return errors.New("DATABASE_URL environment variable not set")
	}

	a.templates = template.New("blank")
	_, err = a.templates.ParseGlob("web/templates/*/*.html")
	if err != nil {
		return
	}
	a.log.Println(a.templates.DefinedTemplates())

	a.db, err = db.SetupConnection(context.Background(), databaseUrl)
	if err != nil {
		return
	}

	return nil
}

func (a *App) Run() error {
	r := mux.NewRouter()
	r.Use(handlers.CompressHandler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.HandleFunc("/", a.getIndex).Methods("GET")
	r.HandleFunc("/teas", a.getTeas).Methods("GET")

	r.HandleFunc("/review", a.postReview).Methods("POST")

	a.log.Println("Listening on", a.listenAddr)
	err := http.ListenAndServe(a.listenAddr, r)
	a.log.Println(err)
	a.db.Close(context.Background())
	return err
}
