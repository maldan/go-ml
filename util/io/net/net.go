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

type Response struct {
	StatusCode int
	Body       io.ReadCloser
	Error      error
	Url        string
}

func (r *Response) JSON(v any) error {
	// Read
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	return err
}

func (r *Response) Bytes() ([]byte, error) {
	// Read
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return body, err
	}

	return body, nil
}

func (r *Response) Close() {
	r.Body.Close()
}

func Get(url string, opts *RequestOptions) Response {
	return Request(url, "GET", opts)
}

func Post(url string, opts *RequestOptions) Response {
	return Request(url, "POST", opts)
}

func Request(url string, method string, options *RequestOptions) Response {
	response := Response{
		Url: url,
	}

	// Get options
	opts := RequestOptions{}
	if options != nil {
		opts = *options
	}
	if opts.Headers == nil {
		opts.Headers = map[string]string{}
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
	if method == "POST" || method == "PATCH" || method == "PUT" {
		switch opts.Data.(type) {
		case []byte:
			inputData = opts.Data.([]byte)
		default:
			out, err := json.Marshal(opts.Data)
			if err != nil {
				response.Error = err
				return response
			}
			opts.Headers["Content-Type"] = "application/json"
			inputData = out
			break
		}
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
