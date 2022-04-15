package app

import (
	"net/http"
	"strconv"

	"github.com/Quantaly/mltea-reviews/app/db"
)

func (a *App) postReview(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	reviewer := r.Form.Get("name")
	if reviewer == "" {
		http.Error(w, "empty name not allowed", http.StatusBadRequest)
		return
	}

	teaID, err := strconv.Atoi(r.Form.Get("tea"))
	if err != nil {
		http.Error(w, "invalid tea id", http.StatusBadRequest)
		return
	}

	rating, err := strconv.Atoi(r.Form.Get("rating"))
	if err != nil || 1 > rating || rating > 5 {
		http.Error(w, "invalid rating", http.StatusBadRequest)
		return
	}

	comment := r.Form.Get("comment")

	_, err = a.db.Exec(r.Context(), db.StmtInsertReview, reviewer, teaID, rating, comment)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
