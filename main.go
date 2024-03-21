package main

import (
	"bookService/auth"
	"bookService/config"
	"bookService/http"
	"bookService/store"
	"log"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatalf("Can't read config file: %v", err)
	}
	mongoStore, err := store.NewMongoStore(conf)
	if err != nil {
		log.Fatalf("main NewMongoStore err: %v", err)
	}
	atKey, err := auth.GenerateECDSAPrivateKey()
	if err != nil {
		log.Fatalf("main generateECDSAPrivateKey atKey err: %v", err)
	}
	rtKey, err := auth.GenerateECDSAPrivateKey()
	if err != nil {
		log.Fatalf("main generateECDSAPrivateKey rtKey err: %v", err)
	}
	middleware := auth.NewAuthMiddleware(atKey, rtKey, mongoStore)
	if err := http.NewServer(mongoStore, middleware); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)

		return
	}
	http.Wait()
}
