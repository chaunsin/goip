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
	"net"
	"net/http"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestClientIP(t *testing.T) {
	type args struct {
		req          *http.Request
		customHeader []XHeader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientIP(tt.args.req, tt.args.customHeader...); got != tt.want {
				t.Errorf("ClientIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewParse(t *testing.T) {
	type args struct {
		proxies      []string
		customHeader []XHeader
	}
	tests := []struct {
		name    string
		args    args
		want    *Parser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParse(tt.args.proxies, tt.args.customHeader...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseForwardedHeader(t *testing.T) {
	type args struct {
		headerValue string
	}
	tests := []struct {
		name    string
		args    args
		want    []ForwardedHeader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseForwardedHeader(tt.args.headerValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseForwardedHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseForwardedHeader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ClientIP(t *testing.T) {
	type fields struct {
		trustedCIDRs atomic.Pointer[[]*net.IPNet]
		xHeader      []XHeader
	}
	type args struct {
		req          *http.Request
		customHeader []XHeader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				trustedCIDRs: tt.fields.trustedCIDRs,
				xHeader:      tt.fields.xHeader,
			}
			if got := p.ClientIP(tt.args.req, tt.args.customHeader...); got != tt.want {
				t.Errorf("ClientIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_SetTrustedProxies(t *testing.T) {
	type fields struct {
		trustedCIDRs atomic.Pointer[[]*net.IPNet]
		xHeader      []XHeader
	}
	type args struct {
		proxies []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				trustedCIDRs: tt.fields.trustedCIDRs,
				xHeader:      tt.fields.xHeader,
			}
			if err := p.SetTrustedProxies(tt.args.proxies); (err != nil) != tt.wantErr {
				t.Errorf("SetTrustedProxies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_isTrustedProxy(t *testing.T) {
	type fields struct {
		trustedCIDRs atomic.Pointer[[]*net.IPNet]
		xHeader      []XHeader
	}
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				trustedCIDRs: tt.fields.trustedCIDRs,
				xHeader:      tt.fields.xHeader,
			}
			if got := p.isTrustedProxy(tt.args.ip); got != tt.want {
				t.Errorf("isTrustedProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_validateHeader(t *testing.T) {
	type fields struct {
		trustedCIDRs atomic.Pointer[[]*net.IPNet]
		xHeader      []XHeader
	}
	type args struct {
		header string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantClientIP string
		wantValid    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				trustedCIDRs: tt.fields.trustedCIDRs,
				xHeader:      tt.fields.xHeader,
			}
			gotClientIP, gotValid := p.validateHeader(tt.args.header)
			if gotClientIP != tt.wantClientIP {
				t.Errorf("validateHeader() gotClientIP = %v, want %v", gotClientIP, tt.wantClientIP)
			}
			if gotValid != tt.wantValid {
				t.Errorf("validateHeader() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestParser_validateXHeader(t *testing.T) {
	type fields struct {
		trustedCIDRs atomic.Pointer[[]*net.IPNet]
		xHeader      []XHeader
	}
	type args struct {
		header string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantClientIP string
		wantValid    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				trustedCIDRs: tt.fields.trustedCIDRs,
				xHeader:      tt.fields.xHeader,
			}
			gotClientIP, gotValid := p.validateXHeader(tt.args.header)
			if gotClientIP != tt.wantClientIP {
				t.Errorf("validateXHeader() gotClientIP = %v, want %v", gotClientIP, tt.wantClientIP)
			}
			if gotValid != tt.wantValid {
				t.Errorf("validateXHeader() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestSetTrustedProxies(t *testing.T) {
	type args struct {
		proxies []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetTrustedProxies(tt.args.proxies); (err != nil) != tt.wantErr {
				t.Errorf("SetTrustedProxies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseCIDRs(t *testing.T) {
	type args struct {
		proxies []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*net.IPNet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCIDRs(tt.args.proxies)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCIDRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCIDRs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIP(tt.args.ip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
