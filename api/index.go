package handler

import (
	"github.com/savioruz/simeru-scraper/config"
	_ "github.com/savioruz/simeru-scraper/docs"
	"github.com/savioruz/simeru-scraper/pkg/server"
	"net/http"
)

// Handler is a function main for vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()

	conf, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Error loading config", http.StatusInternalServerError)
		return
	}

	s := server.NewFiberServer(conf)
	s.Adaptor().ServeHTTP(w, r)
}
