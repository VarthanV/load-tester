package tester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"
)

// config: config for the tester to function
type config struct {
	// The max connections that will be during the peak
	TargetUsers int
	// Duration to reach peak connection after the starting
	// of the connection
	ReachPeakAfter time.Duration
	// Number of users to start with  when the connection starts
	// defaults to
	UsersToStartWith int

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

	db *gorm.DB
}

type Option func(*config)

// Option fn to configure peak  limit
func WithPeakConfig(targetUsers int, reachPeakAfter time.Duration,
	usersToStartWith int) Option {
	return func(c *config) {
		c.ReachPeakAfter = reachPeakAfter
		c.UsersToStartWith = usersToStartWith
		c.TargetUsers = targetUsers
	}
}

// Option fn to configure requests
func WithRequestConfig(url string, body interface{}, acceptedStatusCodes ...int) Option {
	return func(c *config) {
		c.URL = url
		c.Body = body
		c.SuccessStatusCodes = append(c.SuccessStatusCodes, acceptedStatusCodes...)

	}
}

// Option fn to configure custom headers for the request if needed
func WithHeaders(headers map[string]string) Option {
	return func(c *config) {
		h := http.Header{}
		for k, v := range headers {
			h.Set(k, v)
		}
		c.Headers = h
	}
}

func WithDB(db *gorm.DB) Option {
	return func(c *config) {
		c.db = db
	}
}

type driver struct {
	config
	mu                    sync.Mutex
	httpClient            *http.Client
	marshalledBody        []byte
	usersPerMinute        int
	totalNumberOfRequests int
	responseTimeInSeconds []float64
	requestsSucceeded     int
	requestsFailed        int
}

func New(opts ...Option) (*driver, error) {
	d := &driver{
		mu:                    sync.Mutex{},
		totalNumberOfRequests: 0,
		responseTimeInSeconds: make([]float64, 0),
		requestsSucceeded:     0,
		requestsFailed:        0,
	}
	c := config{
		SuccessStatusCodes: []int{http.StatusOK},
		Headers:            http.Header{},
	}

	for _, op := range opts {
		op(&c)
	}

	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConns:      c.TargetUsers,
		MaxConnsPerHost:   c.TargetUsers,
		IdleConnTimeout:   c.ReachPeakAfter,
	}

	// Create an HTTP client using the custom Transport
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Timeout for the HTTP client itself
	}

	d.httpClient = client

	if c.ReachPeakAfter.Minutes() > 0 {
		d.usersPerMinute = (c.TargetUsers - c.UsersToStartWith) / int(c.ReachPeakAfter.Minutes())
	} else {
		d.usersPerMinute = c.TargetUsers
	}

	if c.UsersToStartWith == 0 {
		c.UsersToStartWith = 1
	}

	if c.Body != nil {
		marshalled, err := json.Marshal(c.Body)
		if err != nil {
			log.Println("unable to marshal body ", err)
			return nil, err
		}

		d.marshalledBody = marshalled

	}
	d.config = c

	return d, nil
}

func (d *driver) Run(ctx context.Context) {
	var (
		wg      sync.WaitGroup
		ramupWg sync.WaitGroup
	)

	jobQueue := make(chan struct{}, d.TargetUsers)

	workerCount := d.TargetUsers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobQueue {
				d.doRequestAndReturnStatsDriver(ctx)
			}
		}()
	}

	ramupWg.Add(1)
	go func() {
		defer ramupWg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		usersAdded := 0
		usersToAddPerSecond := d.usersPerMinute / 60
		if d.usersPerMinute%60 != 0 {
			usersToAddPerSecond += 1
		}

		for {
			select {
			case <-ticker.C:
				for i := 0; i < usersToAddPerSecond; i++ {
					if usersAdded >= (d.TargetUsers - d.UsersToStartWith) {
						return
					}
					jobQueue <- struct{}{}
					usersAdded++
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start with initial users
	for i := 0; i < d.UsersToStartWith; i++ {
		jobQueue <- struct{}{}
	}

	ramupWg.Wait()
	close(jobQueue)
	wg.Wait()

	log.Println("Total requests:", d.totalNumberOfRequests)
	r := d.computeReport()
	log.Printf("Report: %+v", r)
}

func (d *driver) doRequestAndReturnStats(ctx context.Context,
	method string, url string, body []byte) (*RequestStat, error) {

	log.Printf("Making request %s %s \n ", d.URL, d.Method)
	stat := RequestStat{}
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx,
		method, url,
		bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("error in creating request %s \n", err.Error())

	}
	res, err := d.httpClient.Do(req)
	if err != nil {
		log.Println("error in doing request", err)
		return nil, err
	}

	log.Println("Response status code is ", res.StatusCode)

	defer res.Body.Close()

	if slices.Contains(d.SuccessStatusCodes, res.StatusCode) {
		stat.IsSuccess = true
	}

	elapsed := time.Since(start)

	stat.TimeTakenInSeconds = elapsed.Seconds()

	return &stat, nil
}

// Given a stat for a request modify the struct variables
func (d *driver) processStat(s *RequestStat) {
	// FIXME: high contention is there in this can be simplified by using atomic
	d.mu.Lock()
	d.totalNumberOfRequests += 1
	d.responseTimeInSeconds = append(d.responseTimeInSeconds, s.TimeTakenInSeconds)

	if s.IsSuccess {
		d.requestsSucceeded += 1
	} else {
		d.requestsFailed += 1
	}

	d.mu.Unlock()

}

func (d *driver) doRequestAndReturnStatsDriver(ctx context.Context) {
	stat, err := d.doRequestAndReturnStats(ctx, d.Method, d.URL, d.marshalledBody)
	if err != nil {
		log.Println("error in doing request ", err)
		d.processStat(&RequestStat{
			IsSuccess: false,
		})
		return
	}
	d.processStat(stat)
}

// Computes report post the load testing is done
func (d *driver) computeReport() *Report {
	r := Report{}

	totalRequests := d.totalNumberOfRequests
	if totalRequests == 0 {
		log.Println("No requests made. Cannot compute report.")
		return &r
	}

	// Compute average response time
	sum := 0.0
	for _, t := range d.responseTimeInSeconds {
		sum += t
	}
	r.AverageResponseTime = sum / float64(totalRequests)

	// Compute peak response time
	r.PeakResponseTime = max(d.responseTimeInSeconds)

	// Compute error rate
	r.ErrorRate = float64(d.requestsFailed) / float64(totalRequests)

	// Compute throughput
	r.Throughput = float64(d.requestsSucceeded) / d.ReachPeakAfter.Seconds()

	// Compute percentiles
	sort.Float64s(d.responseTimeInSeconds)
	r.P50Percentile = percentile(d.responseTimeInSeconds, 50)
	r.P90Percentile = percentile(d.responseTimeInSeconds, 90)
	r.P99Percentile = percentile(d.responseTimeInSeconds, 99)

	return &r
}
