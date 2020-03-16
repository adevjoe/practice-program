// rfc2616
package main

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

type Request struct {
	Method      string
	URI         string
	HTTPVersion string
	Host        string
	RemoteHost  string
	Header      map[string]string
	Body        string
}

type Response struct {
	HTTPVersion  string
	Status       StatusCode
	ReasonPhrase string
	Header       map[string]string
	Body         []byte
}

func NewResponse() *Response {
	return &Response{
		HTTPVersion: HTTP11,
		Header:      map[string]string{"Server": "http-index"},
	}
}

func (res *Response) SetHeader(key, value string) {
	if res.Header == nil {
		res.Header = map[string]string{}
	}
	res.Header[key] = value
}

func (res *Response) SetStatus(code StatusCode) {
	res.Status = code
	res.ReasonPhrase = code.ToReasonPhrase()
}

func (res *Response) SetBody(body []byte) {
	res.SetHeader("Content-Length", strconv.Itoa(len(body)))
	res.Body = body
}

type (
	StatusCode uint16
)

func (code StatusCode) ToReasonPhrase() string {
	return reasonPhrase[code]
}

func (code StatusCode) ToByte() []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(code))
	return bytes.TrimRight(b, "\x00")
}

func (code StatusCode) ToString() string {
	return strconv.Itoa(int(code))
}

const (
	// method
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"

	// status code
	StatusOK          StatusCode = 200
	StatusNotFound    StatusCode = 404
	StatusServerError StatusCode = 500

	// http version
	HTTP11 = "HTTP/1.1"

	// CRLF
	CRLF = "\r\n"
	// SP
	SP = " "
	// LF
	LF = "\n"

	// content type
	ContentHTML = "text/html"
)

// pair of status
var reasonPhrase = map[StatusCode]string{
	StatusOK:          "OK",
	StatusNotFound:    "Not Found",
	StatusServerError: "Server Error",
}
