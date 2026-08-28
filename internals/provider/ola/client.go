package ola

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	ApiKey     string
	BaseUrl    string
	HttpClient *http.Client
}

func NewClient(apiKey string, directionUrl string) (*Client, error) {

	return &Client{
		ApiKey:  apiKey,
		BaseUrl: directionUrl,
		HttpClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}, nil
}

func (c *Client) Directions(
	context context.Context,
	origin string,
	destination string,
	alternatives bool,
) (*Response, error) {

	req, err := http.NewRequestWithContext(
		context,
		http.MethodPost,
		c.BaseUrl,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create OLA Direction Request!")
		return nil, err
	}

	query := req.URL.Query()
	query.Set("origin", origin)
	query.Set("destination", destination)
	query.Set("alternatives", strconv.FormatBool(alternatives))
	query.Set("api_key", c.ApiKey)

	req.URL.RawQuery = query.Encode()

	res, err := c.HttpClient.Do(req)

	if err != nil {
		log.Fatalf("OLA Directions Request Failed!")
		return nil, err
	}

	defer res.Body.Close()
	// fmt.Println("OLA Directions API Response Status:", res)
	// fmt.Println("Request URL:", req.URL.String())
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("OLA Direction API returned status %d", res.StatusCode)
	}

	var response Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("Decoding OLA Directions response failed! %w", err)
	}
	return &response, nil

}
