package fetcher

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinGeckoFetcher_GetPrice(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]map[string]float64{"bitcoin": {"usd": 70000}}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mock.Close()

	f := &CoinGeckoFetcher{
		Client:  mock.Client(),
		BaseURL: mock.URL,
	}

	price, err := f.GetPrice("bitcoin")
	if err != nil {
		t.Fatal(err)
	}
	if price != 70000 {
		t.Fatalf("expected 70000, got %f", price)
	}
}

func TestCoinGeckoFetcher_RealWorld(t *testing.T) {
	f := NewCoinGeckoFetcher() // real API
	symbols := []string{"bitcoin", "ethereum", "dogecoin"}

	for _, s := range symbols {
		price, err := f.GetPrice(s)
		if err != nil {
			t.Errorf("failed to fetch %s: %v", s, err)
			continue
		}
		if price <= 0 {
			t.Errorf("invalid price for %s: %f", s, price)
		} else {
			t.Logf("%s price: %f USD", s, price)
		}
	}
}