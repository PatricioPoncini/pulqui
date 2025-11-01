package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/PatricioPoncini/dolarcito/pkg/services"
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
	return "/dolar"
}

func (c *DolarCommand) Execute(ctx context.Context, chatID int64, args []string) error {
	data, err := c.dolarService.GetExchangeRates()
	if err != nil {
		return c.sender.SendMessage(
			chatID,
			"⚠️ Error intentando traer las cotizaciones del día. Por favor intenta más tarde.",
		)
	}

	message := c.FormatExchangeRates(data)
	return c.sender.SendMessage(chatID, message, "Markdown")
}

func (c *DolarCommand) FormatExchangeRates(data []services.DolarResponse) string {
	message := "💵 *Cotizaciones del día*\n\n"

	var officialSell, blueSell float64

	for _, d := range data {
		message += fmt.Sprintf("🇺🇸 *USD %s*\n", d.Nombre)
		message += fmt.Sprintf("Compra: `$%.2f`\n", d.Compra)
		message += fmt.Sprintf("Venta: `$%.2f`\n\n", d.Venta)

		switch strings.ToLower(d.Nombre) {
		case "oficial":
			officialSell = d.Venta
		case "blue":
			blueSell = d.Venta
		}
	}

	if officialSell > 0 && blueSell > 0 {
		gap := ((blueSell - officialSell) / officialSell) * 100
		message += fmt.Sprintf("💡 *Brecha Blue / Oficial:* `%.2f%%`\n", gap)
	}

	return message
}
