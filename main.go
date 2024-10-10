package main

import (
	"flag"
	"log"
	"read-adviser/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient = telegram.New(tgBotHost, mustToken())

}

func mustToken() string {
	token := flag.String("bot-api-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
