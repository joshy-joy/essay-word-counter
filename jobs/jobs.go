package jobs

import (
	"container/heap"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cenkalti/backoff/v4"
	"github.com/joshy-joy/essay-word-counter/config"
	"github.com/joshy-joy/essay-word-counter/constants"
	"github.com/joshy-joy/essay-word-counter/externals"
	"github.com/joshy-joy/essay-word-counter/utils"
	"github.com/joshy-joy/essay-word-counter/utils/minheap"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	utilsReadFile       = utils.ReadFile
	externalsFetchEssay = externals.FetchEssay
)
var (
	wordFreqMap = make(map[string]int)
	wordFreqMux sync.Mutex
)

func StartWorkerPool(ctx context.Context) error {

	urls, err := utilsReadFile(config.Get().DefaultFilePath)
	if err != nil {
		return err
	}

	jobChan := make(chan string, len(urls))
	h := minheap.NewMinHeap()
	heap.Init(h)

	var wg sync.WaitGroup

	// Start scraping workers
	for i := 0; i < config.Get().WebScrapper.Count; i++ {
		wg.Add(1)
		go func() {
			for _, url := range urls {
				scrapper(ctx, url, jobChan, &wg)
			}
			close(jobChan)
		}()
	}

	// Start word processing workers
	for i := 0; i < config.Get().Tokenizer.Count; i++ {
		wg.Add(1)
		go tokenizer(jobChan, &wg, h)
	}

	wg.Wait()

	result := make([]minheap.Heap, h.Len())
	for i := 0; h.Len() > 0; i++ {
		result[i] = heap.Pop(h).(minheap.Heap)
	}

	formatterJson, err := utils.PrettyPrintJSON(result)
	if err != nil {
		return err
	}
	fmt.Println(formatterJson)
	return nil
}

// Worker function to process each URL
func scrapper(ctx context.Context, url string, jobChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	operation := func() error {
		body, err := externalsFetchEssay(ctx, "GET", url)
		if err != nil {
			log.Printf("error getting url response")
			return err
		}
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			log.Printf("Failed to parse page %s: %v", url, err)
			return err
		}

		var content strings.Builder
		doc.Find("body").Each(func(i int, s *goquery.Selection) {
			// Iterate over all the child nodes of the body tag
			content.WriteString(extractText(s))
		})

		jobChan <- content.String()
		return nil
	}

	// Retry on failure with exponential backoff
	err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5))
	if err != nil {
		log.Printf("Failed to scrape %s after retries: %v", url, err)
	}
}

// Function to count words from each post
func tokenizer(jobChan chan string, wg *sync.WaitGroup, h *minheap.MinHeap) {
	defer wg.Done()
	for content := range jobChan {
		words := getWords(content)
		wordFreqMux.Lock()
		for _, word := range words {
			// condition: to filter words with minimum length
			if len(word) >= config.Get().WordMinLength {
				wordFreqMap[word]++
				// Add the current number to the heap
				heap.Push(h, minheap.Heap{Word: word, Count: wordFreqMap[word]})

				// If heap size exceeds 10, remove the smallest element
				if h.Len() > config.Get().ResultLength {
					heap.Pop(h)
				}
			}
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

// Recursive function to extract text from HTML nodes
func extractText(s *goquery.Selection) string {
	var text strings.Builder
	// Loop through each child node
	s.Contents().Each(func(i int, child *goquery.Selection) {
		if goquery.NodeName(child) == "#text" {
			// If it's a text node, append its content
			if strings.TrimSpace(child.Text()) != constants.Empty {
				text.WriteString(strings.TrimSpace(child.Text()) + " ")
			}
		} else {
			// If it's an element node, extract its child nodes recursively
			text.WriteString(extractText(child))
		}
	})
	return text.String()
}
