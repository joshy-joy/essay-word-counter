package jobs

import (
	"container/heap"
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/joshy-joy/essay-word-counter/config"
	"github.com/joshy-joy/essay-word-counter/utils"
	"github.com/joshy-joy/essay-word-counter/utils/minheap"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"sync"
	"testing"
)

const devConfigFilePath = "../resources/dev/config.yml"

func mockUtilsReadFile(code int) {
	utilsReadFile = func(_ string) ([]string, error) {
		switch code {
		case 1:
			return nil, errors.New("error reading text file")
		default:
			return []string{"https://www.engadget.com/2019/08/25/sony-and-yamaha-sc-1-sociable-cart/",
				"https://www.engadget.com/2019/08/24/trump-tries-to-overturn-ruling-stopping-him-from-blocking-twitte/"}, nil
		}
	}
}

func unMockUtilsReadFile() {
	utilsReadFile = utils.ReadFile
}

func mockFetchEssay(code int) {
	externalsFetchEssay = func(_ context.Context, _, _ string) (io.ReadCloser, error) {
		switch code {
		case 1:
			return nil, errors.New("error getting url response")
		default:
			return io.NopCloser(strings.NewReader("<html><body><p>Test content for test content test </p></body></html>")), nil
		}
	}
}

func unMockFetchEssay() {
	utilsReadFile = utils.ReadFile
}

// Test the scrapper function to ensure it processes pages correctly
func TestScrapper(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	ctx := context.Background()
	jobChan := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	mockFetchEssay(0)
	defer unMockFetchEssay()

	go scrapper(ctx, "https://www.engadget.com/2019/08/25/sony-and-yamaha-sc-1-sociable-cart/", jobChan, &wg)
	content := <-jobChan
	wg.Wait()

	assert.Equal(t, "Test content for test content test ", content, "Expected correct content from scraper")
}

// Test scrapper for error handling
func TestScrapperExternalError(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	ctx := context.Background()
	jobChan := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(1)
	mockFetchEssay(1)
	defer unMockFetchEssay()
	go scrapper(ctx, "https://www.engadget.com/2019/08/25/sony-and-yamaha-sc-1-sociable-cart/", jobChan, &wg)
	close(jobChan)
	wg.Wait()
}

// Test tokenizer to ensure it counts words correctly
func TestTokenizer(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	jobChan := make(chan string, 1)
	h := minheap.NewMinHeap()
	heap.Init(h)
	var wg sync.WaitGroup
	wg.Add(1)

	jobChan <- "joshy joy joshy mike joy sun joshy"
	close(jobChan)

	go tokenizer(jobChan, &wg, h)
	wg.Wait()

	assert.Equal(t, 2, h.Len(), "Expected heap length to be 2")
	top := heap.Pop(h).(minheap.Heap)
	assert.Equal(t, "joy", top.Word, "Expected the top word to be 'joy'")
	assert.Equal(t, 2, top.Count, "Expected the count to be 2")
}

// Test getWords to ensure proper word extraction
func TestGetWords(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	content := "Hello, World! This is a test."
	words := getWords(content)
	expected := []string{"hello", "world", "this", "is", "a", "test"}

	assert.Equal(t, expected, words, "Words extracted are incorrect")
}

// Test extractText function to handle HTML content correctly
func TestExtractText(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	html := "<div><p>Hello</p> <span>world!</span></div>"
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	text := extractText(doc.Selection)
	expected := "Hello world! "

	assert.Equal(t, expected, text, "Extracted text is incorrect")
}

// Test StartWorkerPool with an error case
func TestStartWorkerPoolReadFileErrorCase(t *testing.T) {
	_ = config.InitConfig(devConfigFilePath)
	ctx := context.Background()
	mockUtilsReadFile(1)
	defer unMockUtilsReadFile()

	err := StartWorkerPool(ctx)
	assert.NotNil(t, err, "Expected an error from StartWorkerPool due to file read failure")
}
