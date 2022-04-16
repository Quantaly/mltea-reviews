package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

func SetupConnectionPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	statementNames := []string{
		StmtSelectTeaRatings,
		StmtSelectAllTeaRatings,
		StmtSelectReviews,
		StmtSelectReviewsPaginated,
		StmtSelectTeas,
		StmtSelectFAQEntries,

		StmtInsertReview,
	}
	statementMap := make(map[string]string, len(statementNames))
	for _, stmt := range statementNames {
		sql, err := os.ReadFile(fmt.Sprintf("sql/%s.sql", stmt))
		if err != nil {
			return nil, err
		}
		statementMap[stmt] = string(sql)
	}

	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		for stmt, sql := range statementMap {
			_, err := c.Prepare(ctx, stmt, sql)
			if err != nil {
				return err
			}
		}
		return nil
	}

	config.MaxConns = 15

	return pgxpool.ConnectConfig(ctx, config)
}
