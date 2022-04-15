package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type TeaRating struct {
	Name        string
	Caffeinated bool
	Rating      float64
	ReviewCount int
}

type Review struct {
	Reviewer       string
	Rating         int
	TeaName        string
	TeaCaffeinated bool
	Comment        string
}

type Tea struct {
	ID          int
	Name        string
	Caffeinated bool
}

type FAQEntry struct {
	Question string
	Answer   string
}

const (
	StmtSelectTeaRatings       = "select_tea_ratings"
	StmtSelectAllTeaRatings    = "select_all_tea_ratings"
	StmtSelectReviews          = "select_reviews"
	StmtSelectReviewsPaginated = "select_reviews_paginated"
	StmtSelectTeas             = "select_teas"
	StmtSelectFAQEntries       = "select_faq_entries"

	StmtInsertReview = "insert_review"
)

func SetupConnection(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	for _, stmt := range []string{
		StmtSelectTeaRatings,
		StmtSelectAllTeaRatings,
		StmtSelectReviews,
		StmtSelectReviewsPaginated,
		StmtSelectTeas,
		StmtSelectFAQEntries,

		StmtInsertReview,
	} {
		sql, err := os.ReadFile(fmt.Sprintf("sql/%s.sql", stmt))
		if err != nil {
			conn.Close(ctx)
			return nil, err
		}
		_, err = conn.Prepare(ctx, stmt, string(sql))
		if err != nil {
			conn.Close(ctx)
			return nil, err
		}
	}

	return conn, err
}
