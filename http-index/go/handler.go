package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
			files, err := ioutil.ReadDir(path)
			if err != nil {
				serverError(response, err)
				return
			}

			response.SetBody([]byte(returnFileList(files, request.URI)))
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

func returnFileList(list []os.FileInfo, uri string) string {
	var (
		dirs  []os.FileInfo
		files []os.FileInfo
		s     string
	)

	for _, info := range list {
		if info.IsDir() {
			dirs = append(dirs, info)
		} else {
			files = append(files, info)
		}
	}

	s += "<html><style>body{margin: 50px}</style><body>"
	// is root
	if uri != "" && uri != "/" {
		s += fmt.Sprintf("<a href=\"%s\">%s</a>", getLastDir(uri), "../")
	}
	s += "</br>"
	for _, dir := range dirs {
		s += fmt.Sprintf("<a href=\"%s\">%s</a></br>", strings.TrimRight(uri, "/")+"/"+dir.Name(),
			dir.Name())
	}
	for _, f := range files {
		s += fmt.Sprintf("<a href=\"%s\">%s</a></br>", strings.TrimRight(uri, "/")+"/"+f.Name(),
			f.Name())
	}
	//s += fmt.Sprintf()
	s += "</body></html>"
	return s
}

// server error
func serverError(response *Response, err error) {
	response.SetStatus(StatusServerError)
	response.SetBody([]byte(err.Error()))
}

func getLastDir(path string) string {
	s := path[:strings.LastIndex(strings.TrimRight(path, "/"), "/")]
	if s == "" {
		s = "/"
	}
	return s
}
