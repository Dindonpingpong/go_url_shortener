package main

import (
	"context"
	"log"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest"
)

func main() {
	ctx := context.Background()

	server, err := rest.InitServer(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
