package pokeapi

import (
	"net/http"
	"time"

	"github.com/nurusanwe/pokedexcli/internal/pokecache"
)

// Client struct -
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
	baseURL    string
}

// NewClient initi-
func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:   pokecache.NewCache(cacheInterval),
		baseURL: "https://pokeapi.co/api/v2",
	}
}
