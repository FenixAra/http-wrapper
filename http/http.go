package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	net_url "net/url"
	"strings"
	"time"

	"github.com/FenixAra/go-prom/prom"
	"github.com/FenixAra/go-util/log"
)

type httpwrapper struct {
	c *Config
	l *log.Logger
}

// Function to make HTTP request. method - HTTP method like GET, POST.
// url - HTTP Request URL. req - Request of HTTP request
// res - Pointer to response object
func (h *httpwrapper) MakeRequest(method, url, name string, req, res interface{}) (int, error) {
	if method == http.MethodGet {
		return h.getRequest(method, url, name, res)
	}

	client := &http.Client{
		Timeout: time.Duration(h.c.timeout) * time.Second,
	}

	var retries int
	for {
		var body []byte
		var err error
		var reqBody io.Reader
		if req != nil {
			switch req.(type) {
			case net_url.Values:
				reqBody = strings.NewReader(req.(net_url.Values).Encode())
			default:
				body, err = json.Marshal(req)
				if err != nil {
					h.l.Errorf("Unable to marshal req: %+v. Err: %+v", req, err)
					return 0, err
				}

				reqBody = bytes.NewBuffer(body)
			}
		}

		s := time.Now()
		request, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			h.l.Errorf("Unable to create new HTTP Req. Err: %+v", err)
			continue
		}

		for k, v := range h.c.headers {
			request.Header.Set(k, v)
		}

		response, err := client.Do(request)
		if err != nil {
			h.l.Errorf("Unable to send HTTP Req. Err: %+v", err)
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			time.Sleep(time.Second * time.Duration(int(math.Pow(h.c.retryFactor, float64(retries)))))
			retries++
			if retries < h.c.retries {
				continue
			}

			return 0, err
		}

		go h.l.Log(&log.Log{
			DependancyName: name,
			DependancyType: log.DependancyTypeHTTP,
			TimeTaken:      time.Since(s).Seconds(),
			Title:          url,
			Message:        url,
		})

		if response.StatusCode >= http.StatusInternalServerError {
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			h.l.Errorf("Response code is greater than 500. Code: %d", response.StatusCode)
			time.Sleep(time.Second * time.Duration(int(math.Pow(h.c.retryFactor, float64(retries)))))
			retries++
			if retries < h.c.retries {
				continue
			}

			return response.StatusCode, err
		}

		if response.StatusCode >= http.StatusBadRequest {
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			h.l.Errorf("Response code is between 400 To 499. Code: %d", response.StatusCode)
			return response.StatusCode, err
		}

		prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusSuccess, time.Since(s).Seconds())
		if res != nil {
			content, err := ioutil.ReadAll(response.Body)
			if err != nil {
				h.l.Errorf("Unable to read HTTP Response. Err: %+v", err)
				return response.StatusCode, err
			}

			err = json.Unmarshal(content, &res)
			if err != nil {
				h.l.Errorf("Unable to unmarshal HTTP Response. Err: %+v", err)
				return response.StatusCode, err
			}
		}

		response.Body.Close()
		return response.StatusCode, nil
	}
}

func (h *httpwrapper) getRequest(method, url, name string, res interface{}) (int, error) {
	client := &http.Client{
		Timeout: time.Duration(h.c.timeout) * time.Second,
	}

	var retries int
	for {
		s := time.Now()
		request, err := http.NewRequest(method, url, nil)
		if err != nil {
			h.l.Errorf("Unable to create new HTTP Req. Err: %+v", err)
			return 0, err
		}

		for k, v := range h.c.headers {
			request.Header.Set(k, v)
		}

		response, err := client.Do(request)
		if err != nil {
			h.l.Errorf("Unable to send HTTP Req. Err: %+v", err)
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			time.Sleep(time.Second * time.Duration(int(math.Pow(h.c.retryFactor, float64(retries)))))
			retries++
			if retries < h.c.retries {
				continue
			}

			return 0, err
		}

		go h.l.Log(&log.Log{
			DependancyName: name,
			DependancyType: log.DependancyTypeHTTP,
			TimeTaken:      time.Since(s).Seconds(),
			Title:          url,
			Message:        url,
		})

		if response.StatusCode >= http.StatusInternalServerError {
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			h.l.Errorf("Response code is greater than 500. Code: %d", response.StatusCode)
			time.Sleep(time.Second * time.Duration(int(math.Pow(h.c.retryFactor, float64(retries)))))
			retries++
			if retries < h.c.retries {
				continue
			}

			return response.StatusCode, err
		}

		if response.StatusCode >= http.StatusBadRequest {
			prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusFailed, time.Since(s).Seconds())
			h.l.Errorf("Response code is between 400 To 499. Code: %d", response.StatusCode)
			return response.StatusCode, err
		}

		prom.TrackDependency(prom.DependencyHTTP, name, prom.StatusSuccess, time.Since(s).Seconds())
		if res != nil {
			content, err := ioutil.ReadAll(response.Body)
			if err != nil {
				h.l.Errorf("Unable to read HTTP Response. Err: %+v", err)
				return response.StatusCode, err
			}

			err = json.Unmarshal(content, &res)
			if err != nil {
				h.l.Errorf("Unable to unmarshal HTTP Response. Err: %+v", err)
				return response.StatusCode, err
			}
		}

		response.Body.Close()
		return response.StatusCode, nil
	}
}
