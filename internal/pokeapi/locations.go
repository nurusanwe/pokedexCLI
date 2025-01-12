package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

// struct
// RespShallowLocations -
type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// ListLocations -
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := c.baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// Check cache first
	if cachedData, found := c.cache.Get(url); found {
		var locationsResp RespShallowLocations
		err := json.Unmarshal(cachedData, &locationsResp)
		if err == nil {
			return locationsResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	// Store response in cache
	c.cache.Add(url, dat)

	return locationsResp, nil
}
