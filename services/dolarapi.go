package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type DolarResponse struct {
	Moneda             string    `json:"moneda"`
	Casa               string    `json:"casa"`
	Nombre             string    `json:"nombre"`
	Compra             float64   `json:"compra"`
	Venta              float64   `json:"venta"`
	FechaActualizacion time.Time `json:"fechaActualizacion"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DolarService struct {
	Client  HttpClient
	BaseURL string
}

func NewDolarService(client HttpClient) *DolarService {
	return &DolarService{
		Client:  client,
		BaseURL: os.Getenv("DOLAR_API_URL"),
	}
}

func (s *DolarService) GetExchangeRates() ([]DolarResponse, error) {
	url := fmt.Sprintf("%s/dolares", s.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var data []DolarResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return data, nil
}
