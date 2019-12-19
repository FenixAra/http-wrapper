package main

import (
	"net/http"
	"net/url"

	"github.com/FenixAra/go-log/log"
	http_wrap "github.com/FenixAra/http-wrapper/http"
)

func main() {
	config := log.NewConfig("TestApp")
	l := log.New(config)

	cfg := http_wrap.NewConfig()
	cfg.SetRetries(10)
	cfg.SetTimeout(5)
	cfg.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	wrapper := http_wrap.New(cfg, l)
	req := url.Values{}
	req.Add("id", "test")
	req.Add("name", "name1")

	wrapper.MakeRequest(http.MethodPost, "http://localhost:9000/500", "Google", req, nil)
}
