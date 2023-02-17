package ml_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	net_url "net/url"
)

type RequestOptions struct {
	Data    any
	Headers map[string]string
	Proxy   string
}

type Response[T any] struct {
	StatusCode int
	Body       io.ReadCloser
	Error      error
	Url        string
}

func (r *Response[T]) Unpack() (T, error) {
	out := new(T)

	// Read
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return *out, err
	}

	err = json.Unmarshal(body, out)
	return *out, err
}

func (r *Response[T]) Close() {
	r.Body.Close()
}

func Get[T any](url string, opts RequestOptions) Response[T] {
	return Request[T](url, "GET", opts)
}

func Request[T any](url string, method string, opts RequestOptions) Response[T] {
	response := Response[T]{
		Url: url,
	}

	// Create client
	client := &http.Client{}

	// Set proxy
	if opts.Proxy != "" {
		proxyUrl, _ := net_url.Parse(opts.Proxy)
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Build query
	if method == "GET" || method == "DELETE" {
		mappa := map[string]any{}
		marshal, err := json.Marshal(opts.Data)
		if err != nil {
			response.Error = err
			return response
		}

		err = json.Unmarshal(marshal, &mappa)
		if err != nil {
			response.Error = err
			return response
		}

		response.Url += buildQuery(mappa)
	}

	// Prepare data
	inputData := make([]byte, 0)
	switch opts.Data.(type) {
	case []byte:
		inputData = opts.Data.([]byte)
	default:
		out, err := json.Marshal(opts.Data)
		if err != nil {
			response.Error = err
			return response
		}
		inputData = out
		break
	}

	// Create request
	request, err := http.NewRequest(method, response.Url, bytes.NewBuffer(inputData))
	if err != nil {
		response.Error = err
		return response
	}

	// Fill headers
	for k, v := range opts.Headers {
		request.Header.Set(k, v)
	}

	// Do request
	resp, err := client.Do(request)
	if err != nil {
		response.Error = err
		return response
	}

	// Fill
	response.StatusCode = resp.StatusCode
	response.Body = resp.Body

	return response
}

func buildQuery(data map[string]any) string {
	out := "?"
	for k, v := range data {
		out += fmt.Sprintf("%v=%v&", k, v)
	}
	return out
}
