package main

import (
	"log"
)

func main() {
	log.Println("http server starting...")
	log.Fatal(HttpServerAndListen(":8080"))
}
