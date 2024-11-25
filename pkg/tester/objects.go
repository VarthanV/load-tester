package tester

type Report struct {
	// sum of response time for all requests/total number of requests
	AverageResponseTime float64 `json:"average_response_time"`
	// The longest time taken by the system to respond to a single request during the test period.
	// This metric highlights the worst-case performance scenario.
	PeakResponseTime float64 `json:"peak_response_time"`
	// The percentage of requests that failed during the test, didn't satisfy the
	// success code criteria or error returned during making request
	ErrorRate float64 `json:"error_rate"`
	// The number of requests the system successfully handles per second
	Throughput float64 `json:"throughput"`

	P50Percentile int `json:"p_50_percentile"`

	P90Percentile int `json:"p_90_percentile"`

	P99Percentile int `json:"p_99_percentile"`
}

type RequestStat struct {
	TimeTakenInSeconds float64
	IsSuccess          bool
}
