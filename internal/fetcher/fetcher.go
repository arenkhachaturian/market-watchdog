package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PriceFetcher interface {
	GetPrice(symbol string) (float64, error)
}

type CoinGeckoFetcher struct {
	Client *http.Client
	BaseURL string
}

func NewCoinGeckoFetcher() *CoinGeckoFetcher {
	return &CoinGeckoFetcher{
		Client: &http.Client{Timeout: 10 * time.Second},
		BaseURL: "https://api.coingecko.com/api/v3/simple/price",
	}
}

func (f *CoinGeckoFetcher) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", f.BaseURL, symbol)
	resp, err := f.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, ok := data[symbol]["usd"]
	if !ok {
		return 0, fmt.Errorf("symbol not found")
	}

	return price, nil
}