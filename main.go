package main

import (
	"os"

	"github.com/Quantaly/mltea-reviews/app"
)

func main() {
	a, err := app.New()
	if err == nil {
		a.Run()
	}
	os.Exit(1)
}
