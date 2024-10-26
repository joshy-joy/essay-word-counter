package jobs

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cenkalti/backoff/v4"
	"github.com/joshy-joy/essay-word-counter/config"
	"github.com/joshy-joy/essay-word-counter/utils"
	"github.com/joshy-joy/essay-word-counter/utils/minheap"
)

var (
	wordFreqMap = make(map[string]int)
	wordFreqMux sync.Mutex
)

func InitJobs(ctx context.Context) error {
	urls, err := utils.ReadFile(config.Get().DefaultFilePath)
	if err != nil {
		return err
	}

	jobChan := make(chan string, len(urls))

	h := &minheap.MinHeap{}
	heap.Init(h)

	var wg sync.WaitGroup

	// Start scraping workers
	for i := 0; i < config.Get().WebScrapper.Count; i++ {
		wg.Add(1)
		go func() {
			for _, url := range urls {
				scrpper(ctx, url, jobChan, &wg)
			}
			close(jobChan)
		}()
	}

	// Start word processing workers
	for i := 0; i < config.Get().Tokenizer.Count; i++ {
		wg.Add(1)
		go tokenizer(jobChan, &wg)
	}

	wg.Wait()
	return nil
}

// Worker function to process each URL
func scrpper(ctx context.Context, url string, jobChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	operation := func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}
		client := &http.Client{Timeout: time.Duration(config.Get().External.Timeout)}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to fetch URL %s: %v", url, err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("non-200 status code %d for URL %s", resp.StatusCode, url)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Failed to parse page %s: %v", url, err)
			return err
		}

		content := doc.Find("body").Text()
		jobChan <- content
		return nil
	}

	// Retry on failure with exponential backoff
	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		log.Printf("Failed to scrape %s after retries: %v", url, err)
	}
}

// Function to count words from each post
func tokenizer(jobChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for content := range jobChan {
		words := getWords(content)
		wordFreqMux.Lock()
		for _, word := range words {
			wordFreqMap[word]++
		}
		wordFreqMux.Unlock()
	}
}

// Extract and clean words from the content
func getWords(content string) []string {
	re := regexp.MustCompile(`[^\w\s]+`)
	cleanContent := re.ReplaceAllString(content, "")
	return strings.Fields(strings.ToLower(cleanContent))
}
