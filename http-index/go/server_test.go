package main

import (
	"testing"
)

func TestServer_ParseRequestLine(t *testing.T) {
	var line = []byte("GET /user HTTP/1.1")
	method, uri, version, err := parseRequestLine(line)
	if err != nil {
		t.Error(err.Error())
	}
	if method != MethodGet || uri != "/user" || version != "HTTP/1.1" {
		t.Error()
	}

	var line2 = []byte("POST /user")
	_, _, _, err = parseRequestLine(line2)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error()
	}
}
