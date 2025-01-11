package azhttpclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type AzHttpClient struct {
	client *http.Client
	pat    string
}

func NewAzHttpClient() *AzHttpClient {
	return &AzHttpClient{
		client: &http.Client{},
		pat:    os.Getenv("AZURE_DEVOPS_PAT"),
	}
}

func (c *AzHttpClient) HasValidPat() bool {
	return c.pat != ""
}

func (c *AzHttpClient) SetHeaders(req *http.Request) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(":"+c.pat))

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
}

func Post[TRequest any, TResponse any](c *AzHttpClient, url string, body TRequest) (TResponse, error) {
	var result TResponse

	payload, err := json.Marshal(body)
	if err != nil {
		return result, fmt.Errorf("failed to parse the body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	c.SetHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to perform request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func Get[T any](c *AzHttpClient, url string) (T, error) {
	var result T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to perform request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
