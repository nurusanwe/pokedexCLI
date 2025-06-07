package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nurusanwe/pokedexcli/internal/pokecache"
)

func TestListLocations(t *testing.T) {
	// Mock data
	mockResponse := RespShallowLocations{
		Count:    2,
		Next:     nil,
		Previous: nil,
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "location1", URL: "https://pokeapi.co/api/v2/location-area/1"},
			{Name: "location2", URL: "https://pokeapi.co/api/v2/location-area/2"},
		},
	}
	mockResponseData, _ := json.Marshal(mockResponse)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/location-area" {
			w.WriteHeader(http.StatusOK)
			w.Write(mockResponseData)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create a new Client
	client := NewClient(2*time.Second, 10*time.Second)

	// Override baseURL for testing
	client.baseURL = server.URL

	// Test without cache
	t.Run("Fetch without cache", func(t *testing.T) {
		resp, err := client.ListLocations(nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Validate response
		if resp.Count != len(mockResponse.Results) {
			t.Errorf("expected count %d, got %d", len(mockResponse.Results), resp.Count)
		}
		if resp.Results[0].Name != "location1" {
			t.Errorf("expected first location name 'location1', got %s", resp.Results[0].Name)
		}
	})

	// Test with cache
	t.Run("Fetch with cache", func(t *testing.T) {
		// Second call to ListLocations should hit the cache
		resp, err := client.ListLocations(nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Validate response
		if resp.Count != len(mockResponse.Results) {
			t.Errorf("expected count %d, got %d", len(mockResponse.Results), resp.Count)
		}
	})

	// Test with a specific pageURL
	t.Run("Fetch with specific page URL", func(t *testing.T) {
		pageURL := server.URL + "/location-area"
		resp, err := client.ListLocations(&pageURL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Validate response
		if len(resp.Results) != 2 {
			t.Errorf("expected 2 locations, got %d", len(resp.Results))
		}
	})
}

func TestListLocations_ErrorHandling(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a new Client
	client := NewClient(2*time.Second, 10*time.Second)
	// Override baseURL for testing
	client.baseURL = server.URL

	t.Run("Server error response", func(t *testing.T) {
		_, err := client.ListLocations(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestListExplore(t *testing.T) {
	// Mock response data
	mockResponse := RespLocationsDetail{
		PokemonEncounters: []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
			VersionDetails []struct {
				Version struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"version"`
				MaxChance        int `json:"max_chance"`
				EncounterDetails []struct {
					MinLevel        int   `json:"min_level"`
					MaxLevel        int   `json:"max_level"`
					ConditionValues []any `json:"condition_values"`
					Chance          int   `json:"chance"`
					Method          struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"method"`
				} `json:"encounter_details"`
			} `json:"version_details"`
		}{
			{
				Pokemon: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "pikachu",
					URL:  "https://pokeapi.co/api/v2/pokemon/25/",
				},
				VersionDetails: []struct {
					Version struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"version"`
					MaxChance        int `json:"max_chance"`
					EncounterDetails []struct {
						MinLevel        int   `json:"min_level"`
						MaxLevel        int   `json:"max_level"`
						ConditionValues []any `json:"condition_values"`
						Chance          int   `json:"chance"`
						Method          struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"method"`
					} `json:"encounter_details"`
				}{
					{
						Version: struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						}{
							Name: "red",
							URL:  "https://pokeapi.co/api/v2/version/1/",
						},
						MaxChance: 100,
						EncounterDetails: []struct {
							MinLevel        int   `json:"min_level"`
							MaxLevel        int   `json:"max_level"`
							ConditionValues []any `json:"condition_values"`
							Chance          int   `json:"chance"`
							Method          struct {
								Name string `json:"name"`
								URL  string `json:"url"`
							} `json:"method"`
						}{
							{
								MinLevel:        5,
								MaxLevel:        10,
								ConditionValues: []any{},
								Chance:          50,
								Method: struct {
									Name string `json:"name"`
									URL  string `json:"url"`
								}{
									Name: "walk",
									URL:  "https://pokeapi.co/api/v2/encounter-method/1/",
								},
							},
						},
					},
				},
			},
		},
	}

	// Create a new HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	// Create a new client with a cache
	client := Client{
		httpClient: http.Client{
			Timeout: 5 * time.Second,
		},
		cache: pokecache.NewCache(10 * time.Minute),
	}

	// Override the baseURL with the test server URL
	client.baseURL = ts.URL

	// Call the ListExplore function
	resp, err := client.ListExplore("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response
	if len(resp.PokemonEncounters) != 1 {
		t.Fatalf("expected 1 pokemon encounter, got %d", len(resp.PokemonEncounters))
	}

	if resp.PokemonEncounters[0].Pokemon.Name != "pikachu" {
		t.Fatalf("expected pokemon name to be pikachu, got %s", resp.PokemonEncounters[0].Pokemon.Name)
	}
}

func TestFetchPokemonDetails(t *testing.T) {
	// Mock response data
	mockResponse := PokemonDetails{
		ID:             25,
		Name:           "pikachu",
		BaseExperience: 112,
		Height:         4,
		IsDefault:      true,
		Order:          35,
		Weight:         60,
	}

	// Create a new HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	// Create a new client
	client := Client{
		httpClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
	// Override the baseURL with the test server URL
	client.baseURL = ts.URL

	// Call the FetchPokemonDetails function
	resp, err := client.FetchPokemonDetails("pikachu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response
	if resp.ID != mockResponse.ID {
		t.Fatalf("expected ID %d, got %d", mockResponse.ID, resp.ID)
	}
	if resp.Name != mockResponse.Name {
		t.Fatalf("expected Name %s, got %s", mockResponse.Name, resp.Name)
	}
	if resp.BaseExperience != mockResponse.BaseExperience {
		t.Fatalf("expected BaseExperience %d, got %d", mockResponse.BaseExperience, resp.BaseExperience)
	}
	if resp.Height != mockResponse.Height {
		t.Fatalf("expected Height %d, got %d", mockResponse.Height, resp.Height)
	}
	if resp.IsDefault != mockResponse.IsDefault {
		t.Fatalf("expected IsDefault %v, got %v", mockResponse.IsDefault, resp.IsDefault)
	}
	if resp.Order != mockResponse.Order {
		t.Fatalf("expected Order %d, got %d", mockResponse.Order, resp.Order)
	}
	if resp.Weight != mockResponse.Weight {
		t.Fatalf("expected Weight %d, got %d", mockResponse.Weight, resp.Weight)
	}
}
