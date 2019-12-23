// rfc2616
package main

type Request struct {
	Method      Method
	Path        string
	HTTPVersion string
	Host        string
	Header      map[string]string
	Body        string
}

type Response struct {
	HTTPVersion  string
	Status       StatusCode
	ReasonPhrase string
	Header       map[string]string
	Body         string
}

type (
	Method     string
	StatusCode uint16
)

func (code StatusCode) ToReasonPhrase() string {
	return reasonPhrase[code]
}

const (
	// method
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"

	// status code
	StatusOK          StatusCode = 200
	StatusNotFound    StatusCode = 404
	StatusServerError StatusCode = 500
)

var reasonPhrase = map[StatusCode]string{
	StatusOK:          "OK",
	StatusNotFound:    "Not Found",
	StatusServerError: "Server Error",
}
