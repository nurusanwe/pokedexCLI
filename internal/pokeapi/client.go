package pokeapi

import (
	"net/http"
	"time"

	"github.com/nurusanwe/pokedexcli/internal/pokecache"
)

// Client wraps the details needed to access the PokéAPI.
// The httpClient field performs HTTP requests.
// The cache stores responses to limit network calls.
// The baseURL holds the root API endpoint.
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
	baseURL    string
}

// NewClient creates a Client with the given request timeout and
// cache expiration duration.
func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:   pokecache.NewCache(cacheInterval),
		baseURL: "https://pokeapi.co/api/v2",
	}
}
