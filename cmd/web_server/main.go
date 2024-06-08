package main

import (
	"flag"
	"log"

	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/interfaces/app/web_server"
	"github.com/khailequang334/social_network/internal/interfaces/app/web_server/web_service"
)

var (
	path = flag.String("config", "config.yml", "config path for this service")
)

func main() {
	flag.Parse()
	conf, err := configs.GetWebConfig(*path)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	webSvc, err := web_service.NewWebService(conf)
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
	}
	server := &web_server.WebServer{
		Service: webSvc,
		Port:    conf.Port,
	}
	server.Run()
}
