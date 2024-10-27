package externals

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joshy-joy/essay-word-counter/config"
)

func FetchEssay(ctx context.Context, method, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: time.Duration(config.Get().External.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", url, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code %d for URL %s", resp.StatusCode, url)
	}

	return resp.Body, nil
}
