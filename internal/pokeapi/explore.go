package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type RespLocationsDetail struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
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
	} `json:"pokemon_encounters"`
}

// ListLocations -
func (c *Client) ListExplore(area string) (RespLocationsDetail, error) {
	url := c.baseURL + "/location-area" + "/" + area

	// Check cache first
	if cachedData, found := c.cache.Get(url); found {
		var locationsDetailResp RespLocationsDetail
		err := json.Unmarshal(cachedData, &locationsDetailResp)
		if err == nil {
			return locationsDetailResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespLocationsDetail{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespLocationsDetail{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespLocationsDetail{}, err
	}

	locationsDetailResp := RespLocationsDetail{}
	err = json.Unmarshal(dat, &locationsDetailResp)
	if err != nil {
		return RespLocationsDetail{}, err
	}

	// Store response in cache
	c.cache.Add(url, dat)

	return locationsDetailResp, nil
}
