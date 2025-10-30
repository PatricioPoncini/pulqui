package cron

import (
	"context"
	"log"
	"time"

	"github.com/PatricioPoncini/pulqui/internal/commands"
	"github.com/PatricioPoncini/pulqui/internal/database"
	"github.com/PatricioPoncini/pulqui/pkg/services"
	c "github.com/robfig/cron/v3"
)

var cron *c.Cron

func InitCron(sender commands.MessageSender, dolarService *services.DolarService) error {
	cron = c.New(c.WithSeconds())

	_, err := cron.AddFunc("0 0 17 * * 1-5", func() {
		log.Println("Running daily job:", time.Now().Format("15:04:05"))

		chats, err := database.GetChats(context.Background())
		if err != nil {
			log.Println("Error trying to get chats from db:", err)
			return
		}

		dollarCmd := commands.NewDolarCommand(sender, dolarService)

		data, err := dolarService.GetExchangeRates()
		if err != nil {
			log.Println("Error obtaining exchange rates:", err)
			return
		}

		message := dollarCmd.FormatExchangeRates(data)

		for _, chat := range chats {
			if err := sender.SendMessage(chat.ChatId, message, "Markdown"); err != nil {
				log.Printf("⚠️ Error enviando a chat %v: %v\n", chat.ChatId, err)
			}
		}

		log.Printf("Notifications sent to %d chats\n", len(chats))
	})
	if err != nil {
		log.Fatalf("error running cron job: %v", err)
		return err
	}

	cron.Start()
	log.Println("Cron started successfully")
	return nil
}

func StopCron() {
	cron.Stop()
	log.Println("Cron stopped successfully")
}
