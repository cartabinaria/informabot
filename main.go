package main

import (
	"log"
	"os"

	"github.com/cartabinaria/informabot/bot"
)

const tokenKey = "TOKEN"

func main() {
	token, found := os.LookupEnv(tokenKey)
	if !found {
		log.Fatalf("token not found. please set the %s environment variable",
			tokenKey)
	}

	bot.StartInformaBot(token, false)
}
