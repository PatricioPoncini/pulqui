package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		message += fmt.Sprintf("Venta: `$%.2f`\n", d.Venta)

		now := time.Now()
		if d.FechaActualizacion.Format("2006-01-02") == now.Format("2006-01-02") {
			message += fmt.Sprintf("_(Actualizado hoy a las %s)\n\n_", d.FechaActualizacion.Format("15:04"))
		} else if d.FechaActualizacion.Format("2006-01-02") == now.AddDate(0, 0, -1).Format("2006-01-02") {
			message += fmt.Sprintf("_(Actualizado ayer a las %s)\n\n_", d.FechaActualizacion.Format("15:04"))
		} else {
			message += fmt.Sprintf("_(Actualizado el %s a las %s)\n\n_",
				d.FechaActualizacion.Format("02/01/2006"), d.FechaActualizacion.Format("15:04"))
		}

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
