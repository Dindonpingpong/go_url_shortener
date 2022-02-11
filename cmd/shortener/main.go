package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/filestorage"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewDefaultConfiguration()

	cfg.ParseFlags()
	
	log.Print(cfg.ServerConfig.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := filestorage.NewStorage(cfg.StorageConfig)
	
	if err != nil {
		log.Fatal(err)
	}

	server, err := rest.InitServer(ctx, cfg, storage)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Print("Server Started")

	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
		
		err := storage.PersistStorage()
		
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Print("Server Exited Properly")
}
