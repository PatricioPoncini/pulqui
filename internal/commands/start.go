package commands

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StartCommand struct {
	sender MessageSender
	db     *pgxpool.Pool
}

func NewStartCommand(sender MessageSender, db *pgxpool.Pool) *StartCommand {
	return &StartCommand{sender: sender, db: db}
}

func (c *StartCommand) Name() string {
	return "/start"
}

func (c *StartCommand) Execute(ctx context.Context, chatID int64, args []string) error {
	message := `👋 ¡Hola! Bienvenido al bot de Pulqui

	📊 Para obtener las cotizaciones del dólar del día de hoy, usa el comando:
	
	/dolar
	
	¿Necesitas ayuda? Usa /help para ver todos los comandos disponibles.`

	_, err := c.db.Query(ctx, "INSERT INTO chats (chat_id) VALUES ($1) ON CONFLICT DO NOTHING", strconv.Itoa(int(chatID)))
	if err != nil {
		return err
	}

	return c.sender.SendMessage(chatID, message)
}
