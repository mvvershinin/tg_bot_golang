package main

import (
	"flag"
	"log"
	telegramClient "tg_bot_golang/clients/telegram"
	eventConsumer "tg_bot_golang/consumer/event-consumer"
	"tg_bot_golang/events/telegram"
	"tg_bot_golang/storage/files"
)

const (
	TgBotHost   = "api.telegram.org"
	StoragePath = "file_storage"
	BatchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		telegramClient.New(TgBotHost, mustToken()),
		files.New(StoragePath),
	)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service crashed")
	}

}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for access tg bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token empty")
	}
	return *token
}
