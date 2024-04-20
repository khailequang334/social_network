package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/app/web_app"
	"github.com/khailequang334/social_network/internal/app/web_app/web_service"
)

var (
	path = flag.String("config", "config.yml", "config path for this service")
)

func main() {
	flag.Parse()
	conf, err := configs.GetWebConfig(*path)
	fmt.Println(conf)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	webSvc, err := web_service.NewWebService(conf)
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
	}
	web_app.WebController{
		WebService: *webSvc,
		Port:       conf.Port,
	}.Run()
}
