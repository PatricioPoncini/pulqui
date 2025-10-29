package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PatricioPoncini/pulqui/config"
	"github.com/PatricioPoncini/pulqui/internal/bot"
	"github.com/PatricioPoncini/pulqui/internal/commands"
	"github.com/PatricioPoncini/pulqui/internal/database"
	"github.com/PatricioPoncini/pulqui/internal/telegram"
	"github.com/PatricioPoncini/pulqui/pkg/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Error trying to connect to database: %v", err)
	}
	defer database.Close()

	db := database.GetPool()

	telegramClient := telegram.NewClient(cfg.TelegramToken)

	httpClient := &http.Client{}
	dolarService := services.NewDolarService(httpClient)

	registry := commands.NewRegistry()
	sender := commands.NewTelegramSender(telegramClient)

	registry.Register(commands.NewStartCommand(sender, db))
	registry.Register(commands.NewHelpCommand(sender))
	registry.Register(commands.NewDolarCommand(sender, dolarService))

	botInstance := bot.New(telegramClient, registry)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nInterrupt signal received, stopping bot...")

		done := make(chan struct{})
		go func() {
			cancel()
			db.Close()
			close(done)
		}()

		select {
		case <-done:
			log.Println("Bot stopped")
			os.Exit(0)
		case <-time.After(2 * time.Second):
			log.Println("Timeout reached, forcing shutdown.")
			os.Exit(1)
		}
	}()

	if err := botInstance.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Fatal("Error executing command:", err)
		}
	}
}
