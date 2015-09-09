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

var manager *vpanel.Manager
var monitor *vpanel.HostMonitor

func main() {
	manager = vpanel.NewManager()
	monitor = vpanel.NewHostMonitor(manager)
	monitor.Start()
	defer monitor.Stop()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", Logger(router)))
}
