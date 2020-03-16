package main

import (
	"log"
	"time"
)

func handleConnection(s *Server) {
	for {
		startTime := time.Now()
		request, err := s.Receive()

		// handler request
		response := NewResponse()

		if err != nil {
			serverError(response, err)
		} else {
			handleFile(request, response)
		}

		cost := time.Now().Sub(startTime)
		log.Printf(" |\t %d |\t %s |\t %s | %s\t %s",
			response.Status, cost.String(), request.RemoteHost, request.Method,
			request.URI)
		_ = s.WriteResponse(response)
		s.Close()
		break
	}
}

// TODO handle file
func handleFile(request *Request, response *Response) {
	response.SetStatus(StatusOK)
	response.SetBody([]byte("Hello World!"))
}

// server error
func serverError(response *Response, err error) {
	response.SetStatus(StatusServerError)
	response.SetBody(nil)
}
