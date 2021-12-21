package main

import (
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/internal/app"
)

const(
	address = ":8080"
)

func main() {
	config := app.Config{Address: address}
	db := app.New()

	application := app.App{Config: &config, Db: db}
	application.Start()
}
