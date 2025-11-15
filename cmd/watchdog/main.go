package main

import (
	"log"
	"context"
	"github.com/joho/godotenv"
	inMemoryAlerts "github.com/arenkhachaturian/market-watchdog/internal/store/memory"
	bot "github.com/arenkhachaturian/market-watchdog/internal/bot"
)

func main() {
	myEnv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	botToken := myEnv["TELEGRAM_TOKEN"]
	alertRepo := inMemoryAlerts.NewAlerts()
	tgBot, _ := bot.NewTelegramBot(botToken, alertRepo)
	tgBot.Run(context.Background())
	
}