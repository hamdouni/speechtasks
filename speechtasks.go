// Copyright 2014 Brahim HAMDOUNI. All rights reserved.
// Use of this source code is governed by SIT license
// that can be found in the SIT LICENSE file

// serv is the http backend of speechtasks
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	port  int
)

func browser(url string) error {
	var commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	cmd := exec.Command(run, url)
	return cmd.Start()
}

func init() {
	flag.IntVar(&port, "port", 8080, "port to run the server")
	flag.Parse()
	browser("http://0.0.0.0:" + strconv.Itoa(port) + "/index.html")
}

func loadPage(filename string) ([]byte, string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}
	ext := filepath.Ext(filename)
	mime := mime.TypeByExtension(ext)
	return body, mime, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/"):]
	body, mime, _ := loadPage(filename)
	w.Header().Set("Content-Type", mime)
	fmt.Fprintf(w, "%s", body)
}

func main() {
	println("Server start on port ", port)
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
