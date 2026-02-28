package marketplace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	repoOwner = "Noudea"
	repoName  = "glyph"
	branch    = "main"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

// registry represents the top-level spellbooks/registry.json file.
type registry struct {
	Spellbooks map[string]Spellbook `json:"spellbooks"`
}

// FetchRegistry downloads the registry index in a single HTTP request
// from raw.githubusercontent.com (no API rate limit).
func FetchRegistry() (map[string]Spellbook, error) {
	url := rawURL("spellbooks/registry.json")

	data, err := fetchRaw(url)
	if err != nil {
		return nil, fmt.Errorf("marketplace: fetch registry: %w", err)
	}

	var reg registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("marketplace: parse registry: %w", err)
	}

	return reg.Spellbooks, nil
}

// FetchScript downloads a single script file from a community spellbook.
func FetchScript(id, filename string) ([]byte, error) {
	url := rawURL(fmt.Sprintf("spellbooks/%s/%s", id, filename))

	data, err := fetchRaw(url)
	if err != nil {
		return nil, fmt.Errorf("marketplace: fetch script %s/%s: %w", id, filename, err)
	}
	return data, nil
}

func rawURL(path string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", repoOwner, repoName, branch, path)
}

func fetchRaw(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
