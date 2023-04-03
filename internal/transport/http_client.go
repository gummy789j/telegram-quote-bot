package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HttpRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Cookies map[string]string
	Body    []byte
}

type HttpResponse struct {
	Headers    http.Header
	Cookies    []*http.Cookie
	StatusCode int
	Body       []byte
}

type HttpClient interface {
	Send(ctx context.Context, request *HttpRequest) (*HttpResponse, error)
}

type httpClient struct {
	client *http.Client
}

var _ HttpClient = (*httpClient)(nil)

func defaultTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewHttpClient(proxy ...string) HttpClient {
	transport := defaultTransport()

	if len(proxy) > 0 {
		proxy, _ := url.Parse(proxy[0])
		transport.Proxy = http.ProxyURL(proxy)
	}

	return &httpClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// Use this only when you know the target server is trustable
func NewHttpsClientUnsecure() HttpClient {
	transport := defaultTransport()
	// #nosec
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return &httpClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func GetAPIPath(endpoint, path string) string {
	return fmt.Sprintf("%s%s", endpoint, path)
}

func (c *httpClient) Send(ctx context.Context, request *HttpRequest) (*HttpResponse, error) {
	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for k, v := range request.Params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	for k, v := range request.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("status code: %d", resp.StatusCode)
	}

	response := &HttpResponse{
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		Body:       body,
		StatusCode: resp.StatusCode,
	}

	return response, err
}
