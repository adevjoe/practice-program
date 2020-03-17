package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func handleFile(request *Request, response *Response) {
	// get request path
	path := *Root + request.URI

	// check file
	if f, err := os.Stat(path); os.IsExist(err) || err == nil {
		if f.IsDir() { // list dir
			var s []string
			files, err := ioutil.ReadDir(path)
			if err != nil {
				serverError(response, err)
				return
			}
			for _, info := range files {
				s = append(s, info.Name())
			}

			// TODO list file use html and '<a>' link
			response.SetBody([]byte(fmt.Sprintf("%s", s)))
		} else { // open file
			b, err := ioutil.ReadFile(path)
			if err != nil {
				serverError(response, err)
				return
			}
			response.SetBody(b)
		}
	} else {
		serverError(response, err)
		return
	}
	response.SetStatus(StatusOK)
}

// server error
func serverError(response *Response, err error) {
	response.SetStatus(StatusServerError)
	response.SetBody([]byte(err.Error()))
}
