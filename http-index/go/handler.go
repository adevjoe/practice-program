package main

import "fmt"

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
		request, err := s.Receive()
		// handler request
		response := Response{
			HTTPVersion:  "HTTP/1.1",
			Status:       400,
			ReasonPhrase: "receive request error",
			Header: map[string]string{
				"Server":         "http-index",
				"Content-Type":   "text/html",
				"Content-Length": "",
			},
			Body: err.Error(),
		}
		if err != nil {
			response.Header["content-Length"] = fmt.Sprintf("%d", len(err.Error()))
		}
		_ = s.WriteResponse(response)
		s.Close()
	}
}

func handleFile() {

}
