package http

// Config is used to specify configuration for httpwrapper
type Config struct {
	// HTTP timeout in seconds. Default - 10 second
	timeout int
	// No of Retries in count - Default - 3
	retries int
	// Headers for the http request
	headers     map[string]string
	retryFactor float64
}

// NewConfig Creates and initialises config to default values
func NewConfig() *Config {
	return &Config{
		timeout:     10,
		retries:     3,
		headers:     make(map[string]string),
		retryFactor: 2,
	}
}

// AddHeader is used to add new HTTP header for all requests. k - key of header (Authorisation etc.)
// v - Value of the header
func (c *Config) AddHeader(k, v string) {
	c.headers[k] = v
}

// SetTimeout is used to Set timeout for each HTTP requests
func (c *Config) SetTimeout(timeout int) {
	c.timeout = timeout
}

// SetRetries is used to set number of reties
func (c *Config) SetRetries(retries int) {
	c.retries = retries
}

// SetRetryFactor is used to set retry factor for exponential backoff
func (c *Config) SetRetryFactor(factor float64) {
	c.retryFactor = factor
}
