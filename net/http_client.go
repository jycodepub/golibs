package net

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type HttpClient struct{}

type AuthContext struct {
	Token         string
	FromTokenFile bool // if FromTokenFile, the Token will be the file path
}

type Payload struct {
	Body         string
	FromBodyFile bool // if FromBodyFile, the Body will be the file path
}

type Request struct {
	Payload
	AuthContext
	Headers map[string]string
}

type Response struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}

func (r *Response) IsOK() bool {
	return r.Code == 200
}

func (a *AuthContext) GetBearerToken() (string, error) {
	if a.Token == "" {
		return "", errors.New("no token specified")
	}

	var bearToken string

	if a.FromTokenFile {
		bearToken = fmt.Sprintf("Bearer %s", getToken(a.Token))
	} else {
		bearToken = fmt.Sprintf("Bearer %s", a.Token)
	}

	return bearToken, nil
}

func (p *Payload) getBody() (*strings.Reader, error) {
	if p.FromBodyFile && p.Body != "" {
		bytes, err := os.ReadFile(p.Body)
		if err != nil {
			return nil, err
		}
		return strings.NewReader(string(bytes)), nil
	} else {
		return strings.NewReader(p.Body), nil
	}
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Get(url string, r ...Request) (Response, error) {
	if len(r) > 0 {
		return c.SubmitRequest(GET, url, r[0])
	} else {
		return c.SubmitRequest(GET, url, Request{})
	}
}

func (c *HttpClient) Post(url string, r Request) (Response, error) {
	return c.SubmitRequest(POST, url, r)
}

func (c *HttpClient) Put(url string, r Request) (Response, error) {
	return c.SubmitRequest(PUT, url, r)
}

func (c *HttpClient) Delete(url string, r ...Request) (Response, error) {
	if len(r) > 0 {
		return c.SubmitRequest(DELETE, url, r[0])
	} else {
		return c.SubmitRequest(DELETE, url, Request{})
	}
}

func (c *HttpClient) SubmitRequest(method string, url string, req Request) (Response, error) {
	body, err := req.getBody()
	if err != nil {
		return Response{}, err
	}

	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return Response{}, err
	}
	// Set Bearer token
	if req.Token != "" {
		bearerToken, err := req.GetBearerToken()
		if err != nil {
			return Response{}, err
		}
		r.Header.Set("Authentication", bearerToken)
	}
	// Set headers
	if req.Headers != nil {
		for k, v := range req.Headers {
			r.Header.Set(k, v)
		}
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Code: resp.StatusCode,
		Body: getResponse(resp),
	}, nil
}

func getFile(path string) *os.File {
	apath := path
	if !filepath.IsAbs(path) {
		apath, _ = filepath.Abs(path)
	}
	f, err := os.Open(apath)
	if err != nil {
		panic(err)
	}
	return f
}

func getResponse(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return string(body)
}

func getToken(token string) string {
	f := getFile(token)
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}
