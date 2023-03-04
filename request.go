package what3words

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (w *w3w) request(ctx context.Context, url string, out interface{}) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", w.apiKey)

	resp, err := w.http.Do(request)
	if err != nil {
		return fmt.Errorf("sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request for %s returned unexpected status %s", url, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response body into output: %w", err)
	}

	return nil
}
