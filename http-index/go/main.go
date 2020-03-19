package main

import (
	"flag"
	"log"
)

var Root *string

func main() {
	Root = flag.String("root", "./www", "handle file path")
	port := flag.String("port", "8080", "http port")
	flag.Parse()
	log.Println("http server starting...")
	log.Printf("args root: %s", *Root)
	log.Printf("http listen on http://127.0.0.1:%s", *port)
	log.Fatal(HttpServerAndListen(":" + *port))
}
