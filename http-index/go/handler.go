package main

import "fmt"

func handleConnection(s *Server) {
	for {
		request, err := s.Receive()
		fmt.Println(request)
		// handler request
		response := &Response{
			HTTPVersion:  "HTTP/1.1",
			Status:       400,
			ReasonPhrase: "receive request error",
			Header: map[string]string{
				"Server":         "http-index",
				"Content-Type":   "text/html",
				"Content-Length": "",
			},
		}
		if err != nil {
			response.Header["content-Length"] = fmt.Sprintf("%d", len(err.Error()))
		}
		_ = s.WriteResponse(response)
		s.Close()
		break
	}
}

func handleFile() {

}
