package commands

import "context"

type HelpCommand struct {
	sender MessageSender
}

func NewHelpCommand(sender MessageSender) *HelpCommand {
	return &HelpCommand{sender: sender}
}

func (h *HelpCommand) Name() string {
	return "/help"
}

func (h *HelpCommand) Execute(ctx context.Context, chatID int64, args []string) error {
	return h.sender.SendMessage(chatID, "Comandos disponibles:\n/start\n/dolar\n/help")
}
