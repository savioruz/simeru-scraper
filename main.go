package main

import (
	"github.com/savioruz/simeru-scraper/config"
	_ "github.com/savioruz/simeru-scraper/docs"
	"github.com/savioruz/simeru-scraper/pkg/server"
	"log"
)

// @title Simeru Scraper API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	s := server.NewFiberServer(conf)
	s.ServerStart()
}
