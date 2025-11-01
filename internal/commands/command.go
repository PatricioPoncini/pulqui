package commands

import (
	"context"

	"github.com/PatricioPoncini/dolarcito/internal/telegram"
)

type Command interface {
	Name() string
	Execute(ctx context.Context, chatID int64, args []string) error
}

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) Get(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

type MessageSender interface {
	SendMessage(chatID int64, text string, parseMode ...string) error
}

type TelegramSender struct {
	client *telegram.Client
}

func NewTelegramSender(client *telegram.Client) *TelegramSender {
	return &TelegramSender{client: client}
}

func (t *TelegramSender) SendMessage(chatID int64, text string, parseMode ...string) error {
	return t.client.SendMessage(chatID, text, parseMode...)
}
