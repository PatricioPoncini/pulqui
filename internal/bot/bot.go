package bot

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/PatricioPoncini/pulqui/internal/commands"
	"github.com/PatricioPoncini/pulqui/internal/telegram"
)

type Bot struct {
	client   *telegram.Client
	registry *commands.Registry
}

func New(client *telegram.Client, registry *commands.Registry) *Bot {
	return &Bot{
		client:   client,
		registry: registry,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	log.Println("Bot initiated...")
	offset := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("Bot stopped")
			return ctx.Err()
		default:
			updates, err := b.client.GetUpdates(offset)
			if err != nil {
				log.Printf("Error obtaining updates: %v", err)
				time.Sleep(time.Second * 3)
				continue
			}

			for _, update := range updates {
				b.handleUpdate(ctx, update)
				offset = update.UpdateId + 1
			}

			if len(updates) == 0 {
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update telegram.Update) {
	text := update.Message.Text
	chatID := int64(update.Message.Chat.Id)

	if text == "" {
		return
	}

	parts := strings.Fields(text)
	commandName := parts[0]
	args := parts[1:]

	cmd, exists := b.registry.Get(commandName)
	if !exists {
		if err := b.client.SendMessage(chatID, "Comando no encontrado, por favor vuelva a intentar..."); err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	if err := cmd.Execute(ctx, chatID, args); err != nil {
		log.Printf("Error executing command %s: %v", commandName, err)
	}
}
