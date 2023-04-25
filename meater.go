package meater

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
)

var logger = logrus.New()

// credentials .
type credentials struct {
	username string
	password string
}

// Client .
type Client struct {
	baseURL         string
	credentials     credentials
	httpClient      *resty.Client
	updateInterval  int
	fetchInterval   int
	isAuthenticated bool
	logLevel        string
}

// Response .
type Response struct {
	Status     string           `json:"status"`
	StatusCode int              `json:"statusCode"`
	Data       json.RawMessage  `json:"data,omitempty"`
	Meta       *json.RawMessage `json:"meta,omitempty"`
	Message    *string          `json:"message,omitempty"`
	Help       *string          `json:"help,omitempty"`
}

// Auth .
type Auth struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

// Devices .
type Devices struct {
	Probes []*Probe `json:"devices,omitempty"`
}

// Option .
type Option func(*Client) error

// BaseURL allows overriding of API client baseURL for testing
func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// Username .
func Username(username string) Option {
	return func(c *Client) error {
		c.credentials.username = username
		return nil
	}
}

// Password .
func Password(password string) Option {
	return func(c *Client) error {
		c.credentials.password = password
		return nil
	}
}

// LogLevel allows overriding of API client baseURL for testing
func LogLevel(logLevel string) Option {
	return func(c *Client) error {
		c.logLevel = logLevel
		return nil
	}
}

// auth .
func (c *Client) auth() bool {

	params := map[string]string{
		"email":    c.credentials.username,
		"password": c.credentials.password,
	}
	reqURL := API_VERSION + apiURLs["API_LOGIN"]
	resp, err := c.execute(reqURL, POST, params, true)
	if err != nil {
		return false
	}

	auth := Auth{}
	json.Unmarshal([]byte(resp), &auth)

	c.httpClient.SetAuthToken(auth.Token)

	return true
}

// parseOptions parses the supplied options functions and returns a configured
// *Client instance
func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// New function creates a New Meater client
func New(opts ...Option) (*Client, error) {

	client := Client{
		updateInterval: 7200,
		fetchInterval:  360,
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{}) // &log.TextFormatter{} | &log.JSONFormatter{}

	switch client.logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
		logger.Warnf("[MEATER] Invalid log level supplied: '%s'", logger.GetLevel())
	}

	httpClient := resty.New()

	if client.baseURL != "" {
		httpClient.SetBaseURL(client.baseURL)
	} else {
		httpClient.SetBaseURL(API_SERVER)
	}
	httpClient.
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"User-Agent":      "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate"})

	client.httpClient = httpClient

	if client.auth() {
		return &client, nil
	}

	return &client, nil
}

// ProbesList .
func (c *Client) GetProbes() []*Probe {
	reqURL := API_VERSION + apiURLs["API_DEVICES"]
	resp, err := c.execute(reqURL, GET, map[string]string{}, false)
	if err != nil {
		panic(err)
	}

	logger.Debugf("[MEATER] PROBES: %+v", string(resp))
	devices := Devices{}
	err = json.Unmarshal(resp, &devices)
	if err != nil {
		panic(err)
	}

	if len(devices.Probes) > 0 {
		for _, p := range devices.Probes {
			p.client = c
		}
	}

	return devices.Probes
}

// GetProbeByID .
func (c *Client) GetProbeByID(pID string) *Probe {
	reqURL := API_VERSION + apiURLs["API_DEVICES"] + "/" + pID
	resp, err := c.execute(reqURL, GET, map[string]string{}, false)
	if err != nil {
		panic(err)
	}

	logger.Debugf("[MEATER] PROBE: %+v", string(resp))
	probe := Probe{client: c}
	err = json.Unmarshal(resp, &probe)
	if err != nil {
		panic(err)
	}

	return &probe
}

// Exec method executes a Client instance with the API URL
// Rate limit:
//
//	Recommended: 2 requests per 60 seconds.
//	Maximum: 60 requests per 60 seconds.
func (c *Client) execute(requestUrl string, method string, params map[string]string, jsonEnc bool) ([]byte, error) {
	defer timeTrack("[MEATER][TIMETRK] Executing HTTP Request", logger)

	var resp *resty.Response
	// GET Requests
	if method == "GET" {
		resp, _ = c.httpClient.
			R().
			SetQueryParams(params).
			Get(requestUrl)
	}

	// POST Requests
	if method == "POST" {
		if jsonEnc {
			// POST > JSON Body
			resp, _ = c.httpClient.R().
				SetBody(params).
				Post(requestUrl)
		} else {
			// POST > Form Data
			resp, _ = c.httpClient.R().
				SetFormData(params).
				Post(requestUrl)
		}
	}

	logger.Debugf("[MEATER] HTTP OUTPUT >> %v\n", string([]byte(resp.Body())))

	respParsed := Response{}
	json.Unmarshal([]byte(resp.Body()), &respParsed)

	if respParsed.StatusCode == 200 {
		c.isAuthenticated = true
		return respParsed.Data, nil
	} else {
		c.isAuthenticated = false
		switch error := respParsed.Status; error {
		case apiErrors["API_ERROR_NOT_FOUND"]:
			logger.Debugf("[MEATER] Not Found")
			return nil, errors.New("not found")
		case apiErrors["API_ERROR_UNAUTHORIZED"]:
			logger.Debugf("[MEATER] Client authentication failed")
			return nil, errors.New("client authentication failed")
		case apiErrors["API_ERROR_BAD_REQUEST"]:
			logger.Debugf("[MEATER] Bad Request")
			return nil, errors.New("bad request")
		case apiErrors["API_ERROR_TOO_MANY_REQUESTS"]:
			logger.Debugf("[MEATER] Too Many Requests")
			return nil, errors.New("too many requests")
		case apiErrors["API_ERROR_INTERNAL_SERVER_ERROR"]:
			logger.Debugf("[MEATER] Internal Server Error")
			return nil, errors.New("internal server error")
		default:
			logger.Debugf("[MEATER] Uknown error")
			return nil, errors.New("uknown error")
		}
	}
}
