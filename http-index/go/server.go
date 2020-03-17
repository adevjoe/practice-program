package main

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
)

type Server struct {
	Conn      net.Conn
	BuffR     *bufio.Reader
	BuffW     *bufio.Writer
	StartTime time.Time
	EndTime   time.Time
}

func NewServer(conn net.Conn) *Server {
	s := &Server{
		Conn:      conn,
		BuffR:     bufio.NewReader(conn),
		BuffW:     bufio.NewWriter(conn),
		StartTime: time.Now(),
	}
	return s
}

func (s *Server) Receive() (*Request, error) {
	// handler simple GET request, for HTTP/1.1

	// read request line
	request := new(Request)
	request.RemoteHost = s.Conn.RemoteAddr().String()
	rl, err := s.ReadLine()
	if err != nil {
		return request, err
	}

	request.Method, request.URI, request.HTTPVersion, err = parseRequestLine(rl)
	if err != nil {
		return request, err
	}

	// TODO read header

	return request, nil
}

func parseRequestLine(line []byte) (method, uri, version string, err error) {
	s := strings.Split(string(line), " ")
	// valid s for request line
	if len(s) != 3 {
		return "", "", "", errors.New("invalid request line")
	}
	return s[0], s[1], s[2], nil
}

func (s *Server) ReadLine() ([]byte, error) {
	var line []byte
	for {
		l, more, err := s.BuffR.ReadLine()
		if err != nil {
			return nil, err
		}

		if line == nil && !more {
			return l, nil
		}
		line = append(line, l...)
		if !more {
			break
		}
	}
	return line, nil
}

/**
Response      = Status-Line               ; Section 6.1
			   *(( general-header        ; Section 4.5
				| response-header        ; Section 6.2
				| entity-header ) CRLF)  ; Section 7.1
			   CRLF
			   [ message-body ]          ; Section 7.2
*/
func (s *Server) WriteResponse(response *Response) error {
	// from rfc2616 6.1
	// write status line
	sl, err := parseStatusLine(response)
	if err != nil {
		return err
	}
	_, err = s.BuffW.WriteString(sl)
	if err != nil {
		return err
	}

	// write response header
	h := parseHeader(response.Header)
	_, err = s.BuffW.WriteString(h)
	if err != nil {
		return err
	}

	// CRLF for body
	_, err = s.BuffW.WriteString(CRLF)
	if err != nil {
		return err
	}

	// append body
	_, err = s.BuffW.Write(response.Body)
	if err != nil {
		return err
	}
	return s.BuffW.Flush()
}

func (s *Server) Close() {
	s.EndTime = time.Now()
	//log.Printf("conn %v cost: %s\n", s.Conn.RemoteAddr(), s.EndTime.Sub(s.StartTime))
	_ = s.Conn.Close()
}

func HttpServerAndListen(addr string) error {
	nl, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := nl.Accept()
		if err != nil {
			//log.Printf("%s connect fail, error: %v\n", conn.RemoteAddr(), err)
			continue
		}
		go handleConnection(NewServer(conn))
	}
}

// Status-Line = HTTP-Version SP Status-Code SP Reason-Phrase CRLF
func parseStatusLine(response *Response) (s string, err error) {
	// check field
	if response.HTTPVersion == "" || response.Status <= 0 || response.ReasonPhrase == "" {
		return "", errors.New("status line error")
	}
	// http version
	s += response.HTTPVersion + SP
	// status code
	s += response.Status.ToString() + SP
	// reason phrase
	s += response.ReasonPhrase + SP + CRLF
	return s, nil
}

// (general-header | response-header | entity-header) CRLF
func parseHeader(header map[string]string) string {
	var (
		s     string
		first = true
	)

	// key1: value1
	// key2: value2
	for key, value := range header {
		if !first {
			s += LF
		}
		s += key + ":" + SP + value
		first = false
	}
	s += CRLF
	return s
}
