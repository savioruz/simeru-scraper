package server

type Server interface {
	ServerStart()
}

type Cron interface {
	Start()
	Stop()
}
