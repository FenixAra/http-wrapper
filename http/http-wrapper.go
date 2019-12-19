//go:generate mockgen -destination=mock/mock_http_wrapper.go -package=mock github.com/FenixAra/go-http/http HttpWrapper
package http

import (
	"github.com/FenixAra/go-log/log"
)

// HttpWrapper is the abstracted interface of http
type HttpWrapper interface {
	MakeRequest(method, url, name string, req, res interface{}) (int, error)
}

// New is used to get new HTTP wrapper object
func New(config *Config, l *log.Logger) HttpWrapper {
	return &httpwrapper{
		c: config,
		l: l,
	}
}
