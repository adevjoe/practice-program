package main

import (
	"flag"
	"log"
)

var Root *string

func main() {
	Root = flag.String("root", "./www", "handle file path")
	flag.Parse()
	log.Println("http server starting...")
	log.Printf("args root: %s", *Root)
	log.Fatal(HttpServerAndListen(":8080", ))
}
