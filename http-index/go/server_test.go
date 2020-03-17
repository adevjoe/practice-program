package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestServer_ParseRequestLine(t *testing.T) {
	var line = []byte("GET /user HTTP/1.1")
	method, uri, version, err := parseRequestLine(line)
	if err != nil {
		t.Error(err.Error())
	}
	if method != MethodGet || uri != "/user" || version != "HTTP/1.1" {
		t.Error()
	}

	var line2 = []byte("POST /user")
	_, _, _, err = parseRequestLine(line2)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error()
	}
}

func TestA(t *testing.T) {
	var a uint16 = 200
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, a)
	fmt.Println(bytes.TrimRight(b, "\x00"))
}

func TestListFile(t *testing.T) {
	dir := "./www"
	if f, err := os.Stat(dir); os.IsExist(err) || err == nil {
		t.Logf("name: %s, isDir: %t", f.Name(), f.IsDir())
		files, err := ioutil.ReadDir(f.Name())
		if err != nil {
			t.Error(err)
		}
		for _, file := range files {
			t.Logf("name: %s, isDir: %t", file.Name(), file.IsDir())
		}
	} else {
		os.Mkdir(dir, os.ModePerm)
	}
}

func TestOpenFile(t *testing.T) {
	if b, err := ioutil.ReadFile("./www/media/zelda.jpg"); os.IsExist(err) || err == nil {
		t.Logf("file size: %d", len(b))
	}
}

func TestGetLastDir(t *testing.T) {
	path := "/media/a"
	path2 := "/media/"
	path3 := "/media"
	t.Log(getLastDir(path))
	t.Log(getLastDir(path2))
	t.Log(getLastDir(path3))
}
