package server

import "net/http"

type Server interface {
	ServerStart()
	Adaptor() http.HandlerFunc
}

type Cron interface {
	Start()
	Stop()
}
