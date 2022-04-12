package app

import "net/http"

type indexData struct {
	TopTeas       []ratedTea
	RecentReviews []review

	NonCaffeinatedTeas []tea
	CaffeinatedTeas    []tea

	FAQ []faqEntry
}

type ratedTea struct {
	Name        string
	Caffeinated bool
	Rating      float64
	RatingCount int
}

type review struct {
	Reviewer       string
	Rating         int
	TeaName        string
	TeaCaffeinated bool
	Comment        string
}

type tea struct {
	Id   int
	Name string
}

type faqEntry struct {
	Question string
	Answer   string
}

func (a *App) getIndex(w http.ResponseWriter, r *http.Request) {
	var data indexData
	tx, err := a.db.Begin(r.Context())
	if err != nil {
		a.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// top teas
	data.TopTeas = make([]ratedTea, 0, 10)
	cursor, _ := tx.Query(r.Context(), stmtSelectTopTeas)
	for cursor.Next() {
		var tea ratedTea
		cursor.Scan(&tea.Name, &tea.Caffeinated, &tea.Rating, &tea.RatingCount)
		data.TopTeas = append(data.TopTeas, tea)
	}
	if cursor.Err() != nil {
		a.Logger.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// recent reviews
	data.RecentReviews = make([]review, 0, 5)
	cursor, _ = tx.Query(r.Context(), stmtSelectRecentReviews)
	for cursor.Next() {
		var recentReview review
		cursor.Scan(&recentReview.Reviewer, &recentReview.Rating, &recentReview.TeaName, &recentReview.TeaCaffeinated, &recentReview.Comment)
		data.RecentReviews = append(data.RecentReviews, recentReview)
	}
	if cursor.Err() != nil {
		a.Logger.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// tea list for the dropdown
	data.NonCaffeinatedTeas = make([]tea, 0)
	data.CaffeinatedTeas = make([]tea, 0)
	cursor, _ = tx.Query(r.Context(), stmtSelectAllTeas)
	for cursor.Next() {
		var listedTea tea
		var caffeinated bool
		cursor.Scan(&listedTea.Id, &listedTea.Name, &caffeinated)
		if caffeinated {
			data.CaffeinatedTeas = append(data.CaffeinatedTeas, listedTea)
		} else {
			data.NonCaffeinatedTeas = append(data.NonCaffeinatedTeas, listedTea)
		}
	}
	if cursor.Err() != nil {
		a.Logger.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	// faq
	data.FAQ = make([]faqEntry, 0)
	cursor, _ = tx.Query(r.Context(), stmtSelectFAQEntries)
	for cursor.Next() {
		var entry faqEntry
		cursor.Scan(&entry.Question, &entry.Answer)
		data.FAQ = append(data.FAQ, entry)
	}
	if cursor.Err() != nil {
		a.Logger.Println(cursor.Err())
		http.Error(w, cursor.Err().Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit(r.Context()) // don't particularly care about this error

	err = a.templates.ExecuteTemplate(w, "index.html", &data)
	if err != nil {
		a.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
