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
	message := `ℹ️ *Comandos disponibles:*

	/start — Te da la bienvenida y te explica cómo funciona el bot.
	/dolar — Muestra las cotizaciones actualizadas del dólar (oficial, blue, MEP, etc.).
	/help — Muestra esta lista de comandos.
	
	*Dato extra* 🕔 
	Todos los días a las *17:00 hs* recibirás un mensaje automático con la cotización actualizada y el cierre del mercado.`
	return h.sender.SendMessage(chatID, message, "Markdown")
}
