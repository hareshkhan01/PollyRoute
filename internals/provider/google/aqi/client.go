package aqi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

type Client struct {
	ApiKey     string
	BaseUrl    string
	HttpClient *http.Client
}

func NewClient(apiKey string, AqiUrl string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("Api Key is Empty!")
	}
	if strings.TrimSpace(AqiUrl) == "" {
		return nil, fmt.Errorf("Aqi URL is Empty!")
	}
	return &Client{
		ApiKey:  apiKey,
		BaseUrl: AqiUrl,
		HttpClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}, nil
}

func (c *Client) Aqi(
	context *context.Context,
	location domain.Coordinate,
) (*AQIResponse, error) {

	requestBody := map[string]interface{}{
		"location": map[string]float64{
			"latitude":  location.Lat,
			"longitude": location.Lng,
		},
		"universalAqi": false,
		"extraComputations": []string{"LOCAL_AQI",
			"POLLUTANT_CONCENTRATION"},
		"languageCode": "en",
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to Marshal AQI Request Body! %w", err)
	}
	req, err := http.NewRequestWithContext(
		*context,
		http.MethodPost,
		c.BaseUrl,
		bytes.NewBuffer(requestBodyJson),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create AQI Request! %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	query := req.URL.Query()
	query.Set("key", c.ApiKey)
	req.URL.RawQuery = query.Encode()

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to Request AQI! %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AQI returned with response %d", res.StatusCode)

	}
	var aqiResponse AQIResponse

	if err := json.NewDecoder(res.Body).Decode(&aqiResponse); err != nil {
		return nil, fmt.Errorf("Failed parse AQI Json Body! %w", err)
	}
	return &aqiResponse, nil
}
