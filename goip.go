// MIT License
//
// Copyright (c) 2014 Manuel Martínez-Almeida
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
	"net"
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
	"sync/atomic"
)

// Forwarded standard RFC-7239 https://www.rfc-editor.org/rfc/rfc7239.txt
const Forwarded = "Forwarded"

type XHeader string

const (
	// XForwardedFor X-Forwarded-For
	XForwardedFor XHeader = "X-Forwarded-For"
	// XRealIP commonly used for Nginx、Apache HTTP Server, etc.
	XRealIP XHeader = "X-Real-IP"
	// XClientIP used by Amazon EC2, Heroku, and others
	XClientIP XHeader = "X-Client-IP"
	// CFConnectingIP used by Cloudflare.
	// - https://developers.cloudflare.com/fundamentals/reference/http-request-headers/#cf-connecting-ip
	CFConnectingIP XHeader = "CF-Connecting-IP"
	// FastlyClientIp Fastly CDN and Firebase hosting header when forward to a cloud function
	FastlyClientIp XHeader = "Fastly-Client-Ip"
	// TrueClientIp Akamai and Cloudflare(enterprise plan only)
	TrueClientIp XHeader = "True-Client-Ip"
	// XAppEngineRemoteAddr Google App Engine(GAE)
	XAppEngineRemoteAddr XHeader = "X-Appengine-Remote-Addr"
	// FlyClientIP when running on Fly.io
	FlyClientIP XHeader = "Fly-Client-IP"
)

var (
	defaultParser   *Parser
	remoteIPHeaders = []string{string(XForwardedFor), string(XRealIP)}
	trustedCIDRs    = []string{
		"0.0.0.0/0",      // 0.0.0.0/0 IPv4
		"127.0.0.1/8",    // localhost IPv4
		"10.0.0.0/8",     // 24-bit block IPv4
		"172.16.0.0/12",  // 20-bit block IPv4
		"192.168.0.0/16", // 16-bit block IPv4
		"169.254.0.0/16", // link local address IPv4 https://www.rfc-editor.org/rfc/rfc3927.txt
		"::/0",           // ::/0 IPv6
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6 https://www.rfc-editor.org/rfc/rfc4193.txt
		"fe80::/10",      // link local address IPv6
	}
	defaultTrustedCIDRs = make([]*net.IPNet, 0, len(trustedCIDRs))
)

func init() {
	defaultTrustedCIDRs, _ = parseCIDRs(trustedCIDRs)
	defaultParser, _ = NewParse(nil)
}

type Parser struct {
	trustedCIDRs atomic.Pointer[[]*net.IPNet]
	xHeader      []XHeader
}

func NewParse(proxies []string, customHeader ...XHeader) (*Parser, error) {
	var parser = Parser{
		xHeader: customHeader,
	}
	if proxies == nil || len(proxies) <= 0 {
		parser.trustedCIDRs.Store(&defaultTrustedCIDRs)
	} else {
		cidrs, err := parseCIDRs(proxies)
		if err != nil {
			return nil, err
		}
		parser.trustedCIDRs.Store(&cidrs)
	}
	return &parser, nil
}

func SetTrustedProxies(proxies []string) error {
	return defaultParser.SetTrustedProxies(proxies)
}

func (p *Parser) SetTrustedProxies(proxies []string) (err error) {
	cidrs, err := parseCIDRs(proxies)
	if err != nil {
		return err
	}
	p.trustedCIDRs.Store(&cidrs)
	return
}

func ClientIP(req *http.Request, customHeader ...XHeader) string {
	return defaultParser.ClientIP(req, customHeader...)
}

// ClientIP get client ip
func (p *Parser) ClientIP(req *http.Request, customHeader ...XHeader) string {
	if req == nil || req.Header == nil {
		return ""
	}

	var (
		getValue = req.Header.Get
		xHeader  = append(p.xHeader, customHeader...)
	)

	for _, h := range xHeader {
		value := getValue(string(h))
		if value != "" {
			return value
		}
	}

	ip, _, _ := net.SplitHostPort(textproto.TrimString(req.RemoteAddr))
	remoteIp := net.ParseIP(ip)
	if remoteIp == nil {
		return ""
	}

	trusted := p.isTrustedProxy(remoteIp)
	if !trusted {
		return remoteIp.String()
	}

	for _, h := range remoteIPHeaders {
		ip, valid := p.validateXHeader(getValue(h))
		if valid {
			return ip
		}
	}

	ip, valid := p.validateHeader(getValue(Forwarded))
	if valid {
		return ip
	}
	return remoteIp.String()
}

func (p *Parser) isTrustedProxy(ip net.IP) bool {
	var trustedCIDRs = p.trustedCIDRs.Load()
	if trustedCIDRs == nil {
		return false
	}
	for _, cidr := range *trustedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func parseCIDRs(proxies []string) ([]*net.IPNet, error) {
	if proxies == nil {
		return nil, nil
	}

	var cidr = make([]*net.IPNet, 0, len(proxies))
	for _, trustedProxy := range proxies {
		if !strings.Contains(trustedProxy, "/") {
			var ip = parseIP(trustedProxy)
			if ip == nil {
				return nil, &net.ParseError{Type: "IP address", Text: trustedProxy}
			}

			switch len(ip) {
			case net.IPv4len:
				trustedProxy += "/32"
			case net.IPv6len:
				trustedProxy += "/128"
			}
		}
		_, cidrNet, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			return nil, err
		}
		cidr = append(cidr, cidrNet)
	}
	return cidr, nil
}

func parseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// return ip in a 4-byte representation
		return ipv4
	}
	// return ip in a 16-byte representation or nil
	return parsedIP
}

// validateHeader will parse Forwarded header and return the trusted client IP address
func (p *Parser) validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}

	items, err := ParseForwardedHeader(header)
	if err != nil {
		return "", false
	}
	for i := len(items) - 1; i >= 0; i-- {
		var ipStr = items[i].For
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// Forwarded is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!p.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}

// validateXHeader will parse X-Forwarded-For/X-Real-IP-For header and return the trusted client IP address
func (p *Parser) validateXHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	var items = strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// X-Forwarded-For/X-Real-IP is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!p.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}

// ForwardedHeader represents the structure of a parsed Forwarded header
type ForwardedHeader struct {
	For   string
	By    string
	Host  string
	Proto string
}

// ParseForwardedHeader parses the Forwarded header value and returns a slice of ForwardedHeader
// e.g.
// - Forwarded: for=192.0.2.43, for=198.51.100.17;by=203.0.113.60;proto=http;host=example.com
// - Forwarded: for="[2001:db8:cafe::17]", for=unknown
func ParseForwardedHeader(headerValue string) ([]ForwardedHeader, error) {
	var result []ForwardedHeader

	// Regular expression to match key-value pairs
	pairRegex := regexp.MustCompile(`(?i)(for|by|host|proto)=("[^"]+"|\[[^\]]+\]|[^;,\s]+)`)

	// Split the header by comma to handle multiple forwarded elements
	elements := strings.Split(headerValue, ",")
	for _, element := range elements {
		var forwarded ForwardedHeader
		// Find all key-value pairs in the element
		pairs := pairRegex.FindAllStringSubmatch(element, -1)
		for _, pair := range pairs {
			key := strings.ToLower(pair[1])
			value := strings.Trim(pair[2], "\"")
			value = strings.Trim(value, "[]")
			// Assign the value to the corresponding field in ForwardedHeader
			switch key {
			case "for":
				forwarded.For = value
			case "by":
				forwarded.By = value
			case "host":
				forwarded.Host = value
			case "proto":
				forwarded.Proto = value
			default:
				return nil, fmt.Errorf("unknown forwarded header key: %s", key)
			}
		}
		result = append(result, forwarded)
	}
	return result, nil
}
