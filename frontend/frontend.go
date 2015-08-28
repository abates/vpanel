package main

import (
	"github.com/abates/vpanel"
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {
	manager := vpanel.NewManager()
	manager.Start()
	defer manager.Stop()
	router := NewRouter(manager)
	log.Fatal(http.ListenAndServe(":8080", Logger(router)))
}
