package net

import (
	"bytes"
	"io"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type HttpClient struct{}

type HttpResponse struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}

func (r *HttpResponse) IsOK() bool {
	return r.Code == 200
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Get(url string, ctx *RequestContext) (HttpResponse, error) {
	return c.SubmitRequest(GET, url, ctx)
}

func (c *HttpClient) Post(url string, ctx *RequestContext) (HttpResponse, error) {
	return c.SubmitRequest(POST, url, ctx)
}

func (c *HttpClient) Put(url string, ctx *RequestContext) (HttpResponse, error) {
	return c.SubmitRequest(PUT, url, ctx)
}

func (c *HttpClient) Delete(url string, ctx *RequestContext) (HttpResponse, error) {
	return c.SubmitRequest(DELETE, url, ctx)
}

func (c *HttpClient) SubmitRequest(method string, url string, ctx *RequestContext) (HttpResponse, error) {
	var r *http.Request
	var err error

	if ctx == nil {
		r, err = http.NewRequest(method, url, nil)
		if err != nil {
			return HttpResponse{}, err
		}
	} else {
		// Set request body
		if ctx.Form != nil {
			r, err = http.NewRequest(method, url, bytes.NewReader([]byte(ctx.Form.Encode())))
			if err != nil {
				return HttpResponse{}, err
			}
		} else {
			body, err := ctx.getBody()
			if err != nil {
				return HttpResponse{}, err
			}

			r, err = http.NewRequest(method, url, body)
			if err != nil {
				return HttpResponse{}, err
			}
		}

		// Set Bearer token
		if ctx.Token != "" {
			bearerToken, err := ctx.getBearerToken()
			if err != nil {
				return HttpResponse{}, err
			}
			r.Header.Set("Authorization", bearerToken)
		}

		// Set headers
		if ctx.Headers != nil {
			for k, v := range ctx.Headers {
				r.Header.Set(k, v)
			}
		}
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return HttpResponse{}, err
	}

	return HttpResponse{
		Code: resp.StatusCode,
		Body: getResponse(resp),
	}, nil
}

func getResponse(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return string(body)
}
