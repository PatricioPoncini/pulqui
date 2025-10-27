package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"pulqui/services"
	"strings"

	"github.com/joho/godotenv"
)

var token string

func main() {
	godotenv.Load()
	offset := 0

	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			handleUpdate(update)
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", os.Getenv("TELEGRAM_TOKEN"), offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func handleUpdate(update Update) {
	text := update.Message.Text
	chatID := update.Message.Chat.Id

	if text == "/start" {
		sendMessage(int64(chatID), "Hello World 👋")
	} else if text == "/dolar_hoy" {
		// TODO: Pasar esto a una función aparte
		service := services.NewDolarService(&http.Client{})

		data, err := service.GetExchangeRates()
		if err != nil {
			fmt.Println("Error:", err)
			sendMessage(int64(chatID), "⚠️ Error intentando traer las cotizaciones del día. Por favor intenta maś tarde.")
			return
		}

		message := "💵 *Cotizaciones del día*\n\n"

		for _, d := range data {
			message += fmt.Sprintf("🇺🇸 *USD %s*\n", d.Nombre)
			message += fmt.Sprintf("Compra: `$%.2f`\n", d.Compra)
			message += fmt.Sprintf("Venta: `$%.2f`\n\n", d.Venta)
		}

		sendMessageMarkdown(int64(chatID), message)
	} else {
		sendMessage(int64(chatID), "Command not found.")
	}
}

func sendMessage(chatID int64, text string) {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := fmt.Sprintf("chat_id=%d&text=%s", chatID, text)

	resp, err := http.Post(
		apiUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(payload),
	)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()
}

func sendMessageMarkdown(chatID int64, text string) {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := fmt.Sprintf("chat_id=%d&text=%s&parse_mode=Markdown", chatID, url.QueryEscape(text))

	resp, err := http.Post(
		apiUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(payload),
	)
	if err != nil {
		log.Println("Error enviando mensaje:", err)
		return
	}
	defer resp.Body.Close()
}
