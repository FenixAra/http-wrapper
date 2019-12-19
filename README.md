# http-wrapper
HTTP Wrapper in golang along with mock

[![GoDoc](https://godoc.org/github.com/FenixAra/http-wrapper/http?status.svg)](https://godoc.org/github.com/FenixAra/http-wrapper/http)
[![Go Report Card](https://goreportcard.com/badge/github.com/FenixAra/http-wrapper/http)](https://goreportcard.com/report/github.com/FenixAra/http-wrapper/http)

To get the latest package: 

```sh
go get -u github.com/FenixAra/http-wrapper/http
```

## Usage
```
package main

import (
	"net/http"
	"net/url"

	http_wrap "github.com/FenixAra/go-http/http"
	"github.com/FenixAra/go-log/log"
)

func main() {
	config := log.NewConfig("TestApp")
	l := log.New(config)

	cfg := http_wrap.NewConfig()
	cfg.SetRetries(10)
	cfg.SetTimeout(5)
	cfg.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	wrapper := http_wrap.New(cfg, l)

	wrapper.MakeRequest(http.MethodGet, "https://www.google.com", "Google", nil, nil)
}
```


