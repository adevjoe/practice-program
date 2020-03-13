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

func (s *Server) WriteResponse(response *Response) error {
	b, err := parseResponse(response)
	if err != nil {
		return err
	}
	_, err = s.BuffW.Write(b)
	if err != nil {
		return err
	}
	_ = s.BuffW.Flush()
	return nil
}

func (s *Server) Close() {
	s.EndTime = time.Now()
	//log.Printf("conn %v cost: %s\n", s.Conn.RemoteAddr(), s.EndTime.Sub(s.StartTime))
	_ = s.Conn.Close()
}

func HttpServerAndListen(addr, path string) error {
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

func parseResponse(response *Response) (b []byte, err error) {
	return []byte("HTTP/1.1 200 OK\r\n"), nil
}
