# goip

[![GoDoc](https://godoc.org/github.com/chaunsin/goip?status.svg)](https://godoc.org/github.com/chaunsin/goip) [![Go Report Card](https://goreportcard.com/badge/github.com/chaunsin/goip)](https://goreportcard.com/report/github.com/chaunsin/goip)

A golang library that get client ip from HTTP request

The [gin web framework](https://github.com/gin-gonic/gin) is a very good framework, with its own integration of
the [ClientIP()](https://github.com/gin-gonic/gin/blob/64ead9e6bd924d431f4dd612349bc5e13300e6fc/context.go#L824) method,
but there will be cases in the actual project without the gin framework, just to use the ClientIP() method code would be
too heavy, so stand on the shoulders of the giants, pull out the relevant code, and do some extensions.

## Features

- X-Real-IP rules are supported
- Follows the rule of X-Forwarded-For
- Follows the rule of RFC-7239 standard Forwarded
- The trusted address is allowed to be configured
- Allows getting ip from custom headers
- Exclude local or private address

## Installation

required golang version 1.21+

```shell
go get github.com/chaunsin/goip
```

## Example

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chaunsin/goip"
)

func Example() {
	var serverAddress = "127.0.0.1:8080"

	myParse, err := goip.NewParse([]string{"127.0.0.1"})
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the global trusted proxy addresses，
		// If not set, the default priority is used
		if err := goip.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Println(err)
		}

		// Use the global default configuration
		ip := goip.ClientIP(r)
		// Use custom configuration
		_ = myParse.ClientIP(r)

		// Get the ip address from X-Appengine-Remote-Addr first.
		// If the IP address cannot be obtained, try to get it from another lower priority
		_ = goip.ClientIP(r, goip.XAppEngineRemoteAddr)
		_ = myParse.ClientIP(r, goip.XAppEngineRemoteAddr)

		// Use the custom configuration
		_ = goip.ClientIP(r, goip.XHeader("X-My-IP"))
		_ = myParse.ClientIP(r)

		fmt.Fprintf(w, "Your IP address is %s", ip)
	})
	go func() {
		log.Println(http.ListenAndServe(serverAddress, nil))
	}()

	// execute http request
	req, err := http.NewRequest("GET", "http://"+serverAddress, nil)
	if err != nil {
		panic(err)
	}
	// simulated proxy server added address information
	req.Header.Set("X-Forwarded-For", "123.123.0,1, 123.123.0.2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	// Output:
	// Your IP address is 123.123.0.2
}
```

## How does it work?

The basic principle is to look up the specified object value from the http request header.

The following process is used to search for an ip address:

1. If the user uses a custom request header(
   e.g. [CF-Connecting-IP](https://developers.cloudflare.com/fundamentals/reference/http-request-headers/#cf-connecting-ip)
   、X-Appengine-Remote-Addr...), the fetch is attempted directly from the request header.
2. Get [RemoteAddr](https://github.com/golang/go/blob/48103d97a84d549b44bc4764df6958f73ba5ee02/src/net/http/request.go#L294) to determine whether the address is in the trusted whitelist, and if so, go to Step 3，Otherwise, go to
   Step 4 return RemoteAddr
3. If RemoteAddr is a trusted address, it attempts to get an ip address from X-Real-IP、X-Forward-For、Forwarded, and
   reverse-searches for the first ip address that is not in the trusted address list
4. return ip

## Thanks

- https://github.com/gin-gonic/gin
