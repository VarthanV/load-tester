package tester

import (
	"net/http"
	"time"
)

// Config: Config for the tester to function
type Config struct {
	// The max connections that will be during the peak
	UsersDuringPeakLimit int
	// Duration to reach peak connection after the starting
	// of the connection
	ReachPeakAfter time.Duration
	// Number of users to start with  when the connection starts
	// defaults to
	UsersToStartWith int
	// Duration to ramp up with the provided RamupUserRate
	RampupEvery time.Duration
	// Rate to ramp up with the user
	RamupUserRate int

	// The URL to make request to
	URL string
	// The Http method to make the request
	// Ref: https://www.iana.org/assignments/http-methods/http-methods.xhtml
	Method string
	// Body to send in the request
	Body interface{}

	// Headers if any
	Headers http.Header

	// Accepted http status success codes defaults to 200
	SuccessStatusCodes []int
}

type Option func(*Config)

// Option fn to configure peak  limit
func WithPeakConfig(usersDuringPeakLimit int, reachPeakAfter time.Duration) Option {
	return func(c *Config) {
		c.ReachPeakAfter = reachPeakAfter
		c.UsersDuringPeakLimit = usersDuringPeakLimit
	}
}

// Option  fn to configure ramping up requests
func WithRampupConfig(rampupUserRate int, rampupEvery time.Duration) Option {
	return func(c *Config) {
		c.RampupEvery = rampupEvery
		c.RamupUserRate = rampupUserRate
	}
}

// Option fn to configure requests
func WithRequestConfig(url string, body interface{}, acceptedStatusCodes ...int) Option {
	return func(c *Config) {
		c.URL = url
		c.Body = body
		c.SuccessStatusCodes = append(c.SuccessStatusCodes, acceptedStatusCodes...)

	}
}

// Option fn to configure custom headers for the request if needed
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		h := http.Header{}
		for k, v := range headers {
			h.Set(k, v)
		}
		c.Headers = h
	}
}

type driver struct {
	httpClient *http.Client
	config     Config
}

func New(opts ...Option) *driver {
	c := Config{
		SuccessStatusCodes: []int{http.StatusOK},
	}

	for _, op := range opts {
		op(&c)
	}

	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConns:      c.UsersDuringPeakLimit,
		MaxConnsPerHost:   c.UsersDuringPeakLimit,
		IdleConnTimeout:   c.ReachPeakAfter,
	}

	// Create an HTTP client using the custom Transport
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Timeout for the HTTP client itself
	}

	return &driver{
		httpClient: client,
		config:     c,
	}
}
