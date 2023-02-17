package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyucelen/walker"
)

func buildRequest(start, fetchCount int) (*http.Request, error) {
	url := fmt.Sprintf("https://api.openbrewerydb.org/breweries?page=%d&per_page=%d", start, fetchCount)
	return http.NewRequest(http.MethodGet, url, http.NoBody)
}

func sink(res *http.Response, stop func()) error {
	var payload []map[string]any
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return err
	}

	if len(payload) == 0 {
		stop()
		return nil
	}

	fmt.Println(payload)

	return nil
}

func main() {
	walker.NewApiWalker(http.DefaultClient, buildRequest, sink).Walk()
}
