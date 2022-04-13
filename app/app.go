package app

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"

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

	a.db, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return
	}
	err = a.setupDb()
	if err != nil {
		a.db.Close(context.Background())
		return
	}

	return
}

// stored procedure names
const (
	stmtSelectTopTeas       = "selectTopTeas"
	stmtSelectRecentReviews = "selectRecentReviews"
	stmtSelectAllTeas       = "selectAllTeas"
	stmtSelectFAQEntries    = "selectFAQEntries"
	stmtInsertReview        = "insertReview"
)

func (a *App) setupDb() (err error) {
	_, err = a.db.Prepare(context.Background(), stmtSelectTopTeas, `
		WITH tea_info AS (
			SELECT tea.id, avg(review.rating) AS rating, count(1) AS rating_count
            FROM tea JOIN review ON tea.id = review.tea_id
            GROUP BY tea.id)
		SELECT tea.name, tea.caffeinated, tea_info.rating, tea_info.rating_count
		FROM tea JOIN tea_info ON tea.id = tea_info.id
		ORDER BY tea_info.rating DESC
		LIMIT 10;
	`)
	if err != nil {
		return
	}

	_, err = a.db.Prepare(context.Background(), stmtSelectRecentReviews, `
		SELECT review.reviewer, review.rating, tea.name, tea.caffeinated, review.comment
		FROM review JOIN tea ON review.tea_id = tea.id
		ORDER BY review.id DESC
		LIMIT 5;
	`)
	if err != nil {
		return
	}

	_, err = a.db.Prepare(context.Background(), stmtSelectAllTeas, `
		SELECT id, name, caffeinated
		FROM tea
		ORDER BY name;
	`)
	if err != nil {
		return
	}

	_, err = a.db.Prepare(context.Background(), stmtSelectFAQEntries, `
		SELECT question, answer
		FROM faq
		ORDER BY ordinal;
	`)
	if err != nil {
		return
	}

	_, err = a.db.Prepare(context.Background(), stmtInsertReview, `
		INSERT INTO review (reviewer, tea_id, rating, comment) VALUES ($1, $2, $3, $4);
	`)
	if err != nil {
		return
	}

	return
}

func (a *App) Run() error {
	r := mux.NewRouter()
	r.Use(handlers.CompressHandler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.HandleFunc("/", a.getIndex).Methods("GET")
	r.HandleFunc("/review", a.postReview).Methods("POST")

	a.log.Println("Listening on", a.listenAddr)
	err := http.ListenAndServe(a.listenAddr, r)
	a.log.Println(err)
	a.db.Close(context.Background())
	return err
}
