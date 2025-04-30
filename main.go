package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DerekBelloni/fem_project/internal/app"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("We are running our app")
	http.HandleFunc("/health", HealthCheck)
	server := &http.Server{
		Addr:         ":3572",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is avaiable")
}
