package commands

import (
	"context"
)

type StartCommand struct {
	sender MessageSender
}

func NewStartCommand(sender MessageSender) *StartCommand {
	return &StartCommand{sender: sender}
}

func (c *StartCommand) Name() string {
	return "/start"
}

func (c *StartCommand) Execute(ctx context.Context, chatID int64, args []string) error {
	message := `👋 ¡Hola! Bienvenido al bot de Pulqui

	📊 Para obtener las cotizaciones del dólar del día de hoy, usa el comando:
	
	/dolar
	
	¿Necesitas ayuda? Usa /help para ver todos los comandos disponibles.`

	return c.sender.SendMessage(chatID, message)
}
