package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
)

type Client struct {
	ApiKey     string
	BaseUrl    string
	HttpClient *http.Client
}

func NewClient(apiKey string, weatherUrl string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("Api Key is Empty!")
	}
	if strings.TrimSpace(weatherUrl) == "" {
		return nil, fmt.Errorf("Weather Url is Empty!")
	}
	return &Client{
		ApiKey:  apiKey,
		BaseUrl: weatherUrl,
		HttpClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}, nil
}

func (c *Client) Weather(
	context context.Context,
	location domain.Coordinate,
) (*WeatherResponse, error) {
	req, err := http.NewRequestWithContext(context, http.MethodGet, c.BaseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Weather Request! %w", err)
	}

	query := req.URL.Query()
	query.Set("key", c.ApiKey)
	query.Set("location.latitude", strconv.FormatFloat(location.Lat, 'f', -1, 64))
	query.Set("location.longitude", strconv.FormatFloat(location.Lng, 'f', -1, 64))

	req.URL.RawQuery = query.Encode()

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Weather Request Failed! %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Weather API respond with status code %d", res.StatusCode)
	}

	var response WeatherResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("Failed to parse the Weather Response! %w", err)
	}

	return &response, nil
}
