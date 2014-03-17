package main

import (
    "log"
    "net/http"
    "flag"
)

var (
    certFile = flag.String("cert", "", "")
    keyFile = flag.String("key", "", "")
    root = flag.String("root", "", "")
)

func main() {
    flag.Parse()

    handler := http.FileServer(http.Dir(*root))
    log.Fatal(http.ListenAndServeTLS(":8080", *certFile, *keyFile, handler))
}
