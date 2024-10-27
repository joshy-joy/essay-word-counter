package externals

import (
	"context"
	"github.com/joshy-joy/essay-word-counter/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const devConfigFilePath = "../resources/dev/config.yml"

// Test FetchEssay with a successful HTTP request (status 200)
func TestFetchEssaySuccess(t *testing.T) {
	config.InitConfig(devConfigFilePath)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is a test response"))
	}))
	defer server.Close()

	ctx := context.Background()
	resp, err := FetchEssay(ctx, "GET", server.URL)
	assert.Nil(t, err, "Expected no error for valid request")
	body, _ := io.ReadAll(resp)
	assert.NotNil(t, string(body), "Response body should match expected content")
	resp.Close()
}

// Test FetchEssay with a non-200 status code
func TestFetchEssayError(t *testing.T) {
	config.InitConfig(devConfigFilePath)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	ctx := context.Background()
	resp, err := FetchEssay(ctx, "GET", server.URL)
	assert.NotNil(t, err, "Expected an error for non-200 status code")
	assert.Nil(t, resp, "Expected nil response for non-200 status code")
	assert.Contains(t, err.Error(), "non-200 status code 404", "Error message should contain the status code")
}

// Test FetchEssay with an invalid URL
func TestFetchEssayErrorInvalidURL(t *testing.T) {
	config.InitConfig(devConfigFilePath)
	ctx := context.Background()
	resp, err := FetchEssay(ctx, "GET", "querty")
	assert.NotNil(t, err, "Expected an error due to invalid URL")
	assert.Nil(t, resp, "Expected nil response for invalid URL")
}

// Test FetchEssay with an empty URL
func TestFetchEssayErrorEmptyURL(t *testing.T) {
	config.InitConfig(devConfigFilePath)
	ctx := context.Background()
	resp, err := FetchEssay(ctx, "GET", "")
	assert.NotNil(t, err, "Expected an error due to empty URL")
	assert.Nil(t, resp, "Expected nil response for empty URL")
}
