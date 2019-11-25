.PHONY: build

all: telnet-chat-go telnet-chat-c http-index-go

check-dir:
	@if [ ! -d "build" ]; then mkdir build; fi

telnet-chat-go: check-dir
	@go build -o build/telnet-chat-go telnet-chat/go/*.go
	@echo "telnet-chat-go build success."

telnet-chat-c: check-dir
	@cc telnet-chat/c/main.c -o build/telnet-chat-c
	@echo "telnet-chat-c build success."

http-index-go: check-dir
	@go build -o build/http-index-go http-index/go/*.go
	@echo "http-index-go build success."