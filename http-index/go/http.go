// rfc2616
package main

type Request struct {
	Method      string
	URI         string
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
	StatusCode uint16
)

func (code StatusCode) ToReasonPhrase() string {
	return reasonPhrase[code]
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
)

var reasonPhrase = map[StatusCode]string{
	StatusOK:          "OK",
	StatusNotFound:    "Not Found",
	StatusServerError: "Server Error",
}
