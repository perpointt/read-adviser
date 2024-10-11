package main

import (
	"context"
	"log"
	tgClient "read-adviser/internal/clients/telegram"
	"read-adviser/internal/config"
	"read-adviser/internal/consumer/event-consumer"
	"read-adviser/internal/events/telegram"
	"read-adviser/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	//s := files.New(fileStoragePath)
	s, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	client := tgClient.New(cfg.TgBotApiHost, cfg.TgBotApiToken)

	eventsProcessor := telegram.New(client, s)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.BatchSize)

	if err := consumer.Start(context.TODO()); err != nil {
		log.Fatal("service is stopped", err)
	}
}
