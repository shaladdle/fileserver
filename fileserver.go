package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	auth "github.com/abbot/go-http-auth"
)

var (
	certFile = flag.String("cert", "", "")
	keyFile  = flag.String("key", "", "")
	root     = flag.String("root", "", "")
	port     = flag.Int("port", -1, "")
)

func secret(user, realm string) string {
	if user == "john" {
		// password is "hello"
		return "b98e16cbc3d01734b264adba7baa3bf9"
	}
	return ""
}

func makeAuthHandler(handler http.Handler) auth.AuthenticatedHandlerFunc {
	return func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		handler.ServeHTTP(w, &r.Request)
	}
}

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

	authenticator := auth.NewDigestAuthenticator("example.com", secret)

	log.Printf("args: root = %s, cert = %s, key = %s", absRoot, absCert, absKey)

	hostport := fmt.Sprintf(":%d", *port)

	authHandler := makeAuthHandler(http.FileServer(http.Dir(absRoot)))
	handler := authenticator.Wrap(authHandler)
	log.Fatal(http.ListenAndServeTLS(hostport, absCert, absKey, handler))
}
