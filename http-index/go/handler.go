package main

import "io"

func handleConnection(s *Server) {
	go func(s *Server) {
		for {
			select {
			case <-s.Done:
				return
			}
		}
	}(s)
	for {
		data, err := s.Receive()
		if err != nil {
			if err == io.EOF {
				s.Close()
				return
			}
			return
		}
		if len(data) > 0 {
			s.Send(`HTTP/1.1 200 OK
Date: Sat, 09 Oct 2010 14:28:02 GMT
Server: http-index
Last-Modified: Tue, 01 Dec 2009 20:18:22 GMT
ETag: "51142bc1-7449-479b075b2891b"
Accept-Ranges: bytes
Content-Length: 12
Content-Type: text/txt

Hello World!
`)
			s.Close()
		}
	}
}
