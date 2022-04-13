package app

import (
	"net/http"
	"strconv"

	"github.com/Quantaly/mltea-reviews/app/db"
)

type reviewsData struct {
	Reviews []db.Review

	Page     int
	PrevPage int
	NextPage int
}

func (a *App) getReviews(w http.ResponseWriter, r *http.Request) {
	var data reviewsData

	pagelen := 20 // TODO maybe parameterize?
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	data.Reviews = make([]db.Review, 0, 20)
	reviews_count := 0
	cursor, _ := a.db.Query(r.Context(), db.StmtSelectReviewsPaginated, pagelen, (page-1)*pagelen)
	for cursor.Next() {
		var review db.Review
		cursor.Scan(&reviews_count, &review.Reviewer, &review.Rating, &review.TeaName, &review.TeaCaffeinated, &review.Comment)
		data.Reviews = append(data.Reviews, review)
	}
	if cursor.Err() != nil {
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	data.Page = page
	if page > 1 {
		data.PrevPage = page - 1
	}
	if reviews_count > page*pagelen {
		data.NextPage = page + 1
	}

	err = a.templates.ExecuteTemplate(w, "reviews.html", &data)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
