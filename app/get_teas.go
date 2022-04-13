package app

import (
	"net/http"

	"github.com/Quantaly/mltea-reviews/app/db"
)

type teaListData struct {
	Teas []db.TeaRating

	// Page     int
	// PrevPage int
	// NextPage int
}

func (a *App) getTeas(w http.ResponseWriter, r *http.Request) {
	var data teaListData

	data.Teas = make([]db.TeaRating, 0)
	cursor, _ := a.db.Query(r.Context(), db.StmtSelectAllTeaRatings)
	for cursor.Next() {
		var rating db.TeaRating
		cursor.Scan(&rating.Name, &rating.Caffeinated, &rating.Rating, &rating.ReviewCount)
		data.Teas = append(data.Teas, rating)
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
