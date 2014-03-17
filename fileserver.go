package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

var (
	certFile = flag.String("cert", "", "")
	keyFile  = flag.String("key", "", "")
	root     = flag.String("root", "", "")
	port     = flag.Int("port", -1, "")
)

func main() {
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}

	absKey, err := filepath.Abs(*keyFile)
	if err != nil {
		log.Fatal(err)
	}

	absCert, err := filepath.Abs(*certFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("args: root = %s, cert = %s, key = %s", absRoot, absCert, absKey)

	handler := http.FileServer(http.Dir(absRoot))
	hostport := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServeTLS(hostport, absCert, absKey, handler))
}
