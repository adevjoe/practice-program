package main

import (
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

func (s *Server) Receive() (string, error) {
	data := make([]byte, 1000)
	_, err := s.Conn.Read(data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) Send(msg string) {
	_, _ = s.Conn.Write([]byte(msg))
}

func (s *Server) Close() {
	s.EndTime = time.Now()
	//log.Printf("conn %v cost: %s\n", s.Conn.RemoteAddr(), s.EndTime.Sub(s.StartTime))
	s.Done <- true
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
