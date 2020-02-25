package main

import (
	"io"
	"net"
	"time"
)

type Server struct {
	Conn      net.Conn
	Done      chan bool
	StartTime time.Time
	EndTime   time.Time
}

func NewServer(conn net.Conn) *Server {
	s := &Server{
		Conn:      conn,
		Done:      make(chan bool),
		StartTime: time.Now(),
	}
	return s
}

func (s *Server) Receive() (request Request, err error) {
	data := make([]byte, 1000)
	_, err = s.Conn.Read(data)
	if err != nil {
		if err == io.EOF {
			s.Close()
			return
		}
		return
	}
	return parseToRequest(data)
}

func (s *Server) WriteResponse(response Response) error {
	b, err := parseResponse(response)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() {
	s.EndTime = time.Now()
	//log.Printf("conn %v cost: %s\n", s.Conn.RemoteAddr(), s.EndTime.Sub(s.StartTime))
	s.Done <- true
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

func parseToRequest(b []byte) (request Request, err error) {

}

func parseResponse(response Response) (b []byte, err error) {

}
