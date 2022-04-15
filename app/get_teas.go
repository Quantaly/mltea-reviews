package app

import (
	"net/http"

	"github.com/Quantaly/mltea-reviews/app/db"
)

type teaListData struct {
	Teas            []db.TeaRating
	UnavailableTeas []db.TeaRating
}

func (a *App) getTeas(w http.ResponseWriter, r *http.Request) {
	var data teaListData

	data.Teas = make([]db.TeaRating, 0)
	data.UnavailableTeas = make([]db.TeaRating, 0)
	cursor, _ := a.db.Query(r.Context(), db.StmtSelectAllTeaRatings)
	for cursor.Next() {
		var rating db.TeaRating
		var available bool
		cursor.Scan(&rating.Name, &rating.Caffeinated, &rating.Rating, &rating.ReviewCount, &available)
		if available {
			data.Teas = append(data.Teas, rating)
		} else {
			data.UnavailableTeas = append(data.UnavailableTeas, rating)
		}
	}
	if cursor.Err() != nil {
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	err := a.templates.ExecuteTemplate(w, "teas.html", &data)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
