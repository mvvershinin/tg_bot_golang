package tg_bot_golang

import (
	"flag"
	"log"

	"tg_bot_golang/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient = telegram.New(tgBotHost, mustToken())

	//fetcher = fetcher.New(tgClient)

	//processor = processor.New(tgClient)

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for access tg bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token empty")
	}
	return "string"
}
