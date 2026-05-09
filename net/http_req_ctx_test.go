package net

import "testing"

func TestRequestContext(t *testing.T) {
	ctx := RequestContext{}
	ctx.AddHeader("User-Agent", "Mozilla/5.0")
	ctx.AddHeader("Accept", "text/html,application/json")
	ctx.AddToken("token")
	ctx.AddPayLoad("payload")

	bearToken, err := ctx.getBearerToken()
	if err != nil {
		t.Fail()
	}
	if bearToken != "Bearer token" {
		t.Fail()
	}

	if ctx.Headers["User-Agent"] != "Mozilla/5.0" {
		t.Fail()
	}

	if ctx.Headers["Accept"] != "text/html,application/json" {
		t.Fail()
	}

	if ctx.Body != "payload" {
		t.Fail()
	}
}
