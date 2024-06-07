// MIT License
//
// Copyright (c) 2024 chaunsin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package goip

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Example() {
	var serverAddress = "127.0.0.1:8080"

	myParse, err := NewParse([]string{"127.0.0.1"})
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the global trusted proxy addresses，
		// If not set, the default priority is used
		if err := SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Println(err)
		}

		// Use the global default configuration
		ip := ClientIP(r)
		// Use custom configuration
		_ = myParse.ClientIP(r)

		// Get the ip address from X-Appengine-Remote-Addr first.
		// If the IP address cannot be obtained, try to get it from another lower priority
		_ = ClientIP(r, XAppEngineRemoteAddr)
		_ = myParse.ClientIP(r, XAppEngineRemoteAddr)

		// Use the custom configuration
		_ = ClientIP(r, XHeader("X-My-IP"))
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
