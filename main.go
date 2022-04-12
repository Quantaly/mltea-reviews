package main

import (
	"github.com/Quantaly/mltea-reviews/app"
)

func main() {
	a, err := app.InitApp()
	if err != nil {
		a.Logger.Fatalln(err)
	}

	a.Logger.Fatalln(a.Run())
}
