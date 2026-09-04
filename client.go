package hypixel

import (
	"net/http"
)

type PreRequestHook func(request Request) (Response, error)
type Callback func(request Request, response Response, err error) (Response, error)

type Client struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	rate           *RateLimit
	preRequestHook PreRequestHook
	callBack       Callback
}

// https://api.hypixel.net/
func NewClient(key string, rate *RateLimit) *Client {
	return &Client{
		baseURL:    "https://api.hypixel.net/v2/",
		apiKey:     key,
		httpClient: http.DefaultClient,
		rate:       rate,
	}
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}

func (c *Client) GetAPIKey() string {
	return c.apiKey
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) GetRate() *RateLimit {
	return c.rate
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

func (c *Client) SetRate(rate *RateLimit) {
	c.rate = rate
}

func (c *Client) SetPreRequestHook(beforeSend PreRequestHook) {
	c.preRequestHook = beforeSend
}

func (c *Client) SetCallback(callBack Callback) {
	c.callBack = callBack
}

func (c *Client) authHeader(header ...http.Header) http.Header {
	var h http.Header
	if len(header) == 0 {
		h = http.Header{}
	} else {
		h = header[0]
	}
	h.Set("API-Key", c.apiKey)
	return h
}
