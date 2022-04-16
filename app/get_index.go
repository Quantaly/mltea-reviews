package app

import (
	"net/http"

	"github.com/Quantaly/mltea-reviews/app/db"
)

type indexData struct {
	TopTeas       []db.TeaRating
	RecentReviews []db.Review

	NonCaffeinatedTeas []db.Tea
	CaffeinatedTeas    []db.Tea

	FAQ []db.FAQEntry
}

func (a *App) getIndex(w http.ResponseWriter, r *http.Request) {
	var data indexData
	tx, err := a.db.Begin(r.Context())
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// top teas
	data.TopTeas = make([]db.TeaRating, 0, 10)
	cursor, _ := tx.Query(r.Context(), db.StmtSelectTeaRatings)
	for cursor.Next() {
		var rating db.TeaRating
		cursor.Scan(&rating.Name, &rating.Caffeinated, &rating.Rating, &rating.ReviewCount)
		data.TopTeas = append(data.TopTeas, rating)
	}
	if cursor.Err() != nil {
		tx.Rollback(r.Context())
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// recent reviews
	data.RecentReviews = make([]db.Review, 0, 5)
	cursor, _ = tx.Query(r.Context(), db.StmtSelectReviews)
	for cursor.Next() {
		var review db.Review
		cursor.Scan(&review.Reviewer, &review.Rating, &review.TeaName, &review.TeaCaffeinated, &review.Comment)
		data.RecentReviews = append(data.RecentReviews, review)
	}
	if cursor.Err() != nil {
		tx.Rollback(r.Context())
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// tea list for the dropdown
	data.NonCaffeinatedTeas = make([]db.Tea, 0)
	data.CaffeinatedTeas = make([]db.Tea, 0)
	cursor, _ = tx.Query(r.Context(), db.StmtSelectTeas)
	for cursor.Next() {
		var tea db.Tea
		cursor.Scan(&tea.ID, &tea.Name, &tea.Caffeinated)
		if tea.Caffeinated {
			data.CaffeinatedTeas = append(data.CaffeinatedTeas, tea)
		} else {
			data.NonCaffeinatedTeas = append(data.NonCaffeinatedTeas, tea)
		}
	}
	if cursor.Err() != nil {
		tx.Rollback(r.Context())
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// faq
	data.FAQ = make([]db.FAQEntry, 0)
	cursor, _ = tx.Query(r.Context(), db.StmtSelectFAQEntries)
	for cursor.Next() {
		var entry db.FAQEntry
		cursor.Scan(&entry.Question, &entry.Answer)
		data.FAQ = append(data.FAQ, entry)
	}
	if cursor.Err() != nil {
		tx.Rollback(r.Context())
		a.log.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		tx.Rollback(r.Context())
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.templates.ExecuteTemplate(w, "index.html", &data)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
