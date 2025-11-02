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
	message := `👋 ¡Hola! Bienvenido a *Dolarcito Bot* 🇦🇷💵

	📊 Con este bot podés consultar las cotizaciones del dólar actualizadas en cualquier momento.
	
	🕔 Además, todos los días a las *17:00 hs* recibirás automáticamente un mensaje con la cotización actualizada y el cierre del mercado.
	
	Comandos disponibles:
	/dolar — Muestra las cotizaciones del dólar del día.
	/help — Explica todos los comandos disponibles.`

	_, err := c.db.Query(ctx, "INSERT INTO chats (chat_id) VALUES ($1) ON CONFLICT DO NOTHING", strconv.Itoa(int(chatID)))
	if err != nil {
		return err
	}

	return c.sender.SendMessage(chatID, message, "Markdown")
}
