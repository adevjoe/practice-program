package main

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
		data := s.Receive()
		if len(data) > 0 {
			print(data)
			s.Send(`HTTP/1.1 200 OK
Date: Sat, 09 Oct 2010 14:28:02 GMT
Server: http-index
Accept-Ranges: bytes
Content-Length: 12
Content-Type: text/txt

Hello World!
`)
			s.Close()
		}
	}
}

func handleFile() {

}
