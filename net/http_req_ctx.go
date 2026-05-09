package net

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type AuthToken struct {
	Token     string
	TokenFile string // token can be in a file
}

type Payload struct {
	Body     string
	BodyFile string
	Form     url.Values
}

type RequestContext struct {
	Payload
	AuthToken
	Headers map[string]string
}

func (ctx *RequestContext) AddHeader(key string, value string) *RequestContext {
	if ctx.Headers == nil {
		ctx.Headers = make(map[string]string)
	}
	ctx.Headers[key] = value
	return ctx
}

func (ctx *RequestContext) AddFormField(key string, value string) *RequestContext {
	if ctx.Form == nil {
		ctx.Form = make(url.Values)
	}
	ctx.Form.Add(key, value)
	return ctx
}

func (ctx *RequestContext) AddPayLoad(payload string) *RequestContext {
	ctx.Payload.Body = payload
	return ctx
}

func (ctx *RequestContext) AddPayLoadFile(filepath string) *RequestContext {
	ctx.Payload.BodyFile = filepath
	return ctx
}

func (ctx *RequestContext) AddToken(token string) *RequestContext {
	ctx.AuthToken.Token = token
	return ctx
}

func (ctx *RequestContext) AddTokenFile(filepath string) *RequestContext {
	ctx.AuthToken.TokenFile = filepath
	return ctx
}

func (a *AuthToken) getBearerToken() (string, error) {
	if a.Token == "" && a.TokenFile == "" {
		return "", errors.New("no token specified")
	}

	var bearToken string

	if a.TokenFile != "" {
		bearToken = fmt.Sprintf("Bearer %s", getTokenFromFile(a.TokenFile))
	} else {
		bearToken = fmt.Sprintf("Bearer %s", a.Token)
	}

	return bearToken, nil
}

func (p *Payload) getBody() (*strings.Reader, error) {
	if p.Body != "" {
		return strings.NewReader(p.Body), nil
	} else if p.BodyFile != "" {

		body, err := os.ReadFile(p.BodyFile)
		if err != nil {
			return nil, err
		}
		return strings.NewReader(string(body)), nil
	}

	return nil, errors.New("no body specified")
}

func getTokenFromFile(filepath string) string {
	f := openFile(filepath)
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func openFile(path string) *os.File {
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
