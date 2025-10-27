package commands

import (
	"context"
	"fmt"

	"github.com/PatricioPoncini/pulqui/pkg/services"
)

type DolarCommand struct {
	sender       MessageSender
	dolarService *services.DolarService
}

func NewDolarCommand(sender MessageSender, dolarService *services.DolarService) *DolarCommand {
	return &DolarCommand{
		sender:       sender,
		dolarService: dolarService,
	}
}

func (c *DolarCommand) Name() string {
	return "/dolar_hoy"
}

func (c *DolarCommand) Execute(ctx context.Context, chatID int64, args []string) error {
	data, err := c.dolarService.GetExchangeRates()
	if err != nil {
		return c.sender.SendMessage(
			chatID,
			"⚠️ Error intentando traer las cotizaciones del día. Por favor intenta más tarde.",
		)
	}

	message := c.formatExchangeRates(data)
	return c.sender.SendMessage(chatID, message, "Markdown")
}

func (c *DolarCommand) formatExchangeRates(data []services.DolarResponse) string {
	message := "💵 *Cotizaciones del día*\n\n"

	for _, d := range data {
		message += fmt.Sprintf("🇺🇸 *USD %s*\n", d.Nombre)
		message += fmt.Sprintf("Compra: `$%.2f`\n", d.Compra)
		message += fmt.Sprintf("Venta: `$%.2f`\n\n", d.Venta)
	}

	return message
}
