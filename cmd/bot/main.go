package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PatricioPoncini/pulqui/config"
	"github.com/PatricioPoncini/pulqui/internal/bot"
	"github.com/PatricioPoncini/pulqui/internal/commands"
	"github.com/PatricioPoncini/pulqui/internal/telegram"
	"github.com/PatricioPoncini/pulqui/pkg/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	telegramClient := telegram.NewClient(cfg.TelegramToken)

	httpClient := &http.Client{}
	dolarService := services.NewDolarService(httpClient)

	registry := commands.NewRegistry()
	sender := commands.NewTelegramSender(telegramClient)

	registry.Register(commands.NewStartCommand(sender))
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
		cancel()
	}()

	if err := botInstance.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Fatal("Error executing command:", err)
		}
	}
}
