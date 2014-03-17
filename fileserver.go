package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	auth "github.com/abbot/go-http-auth"
)

var (
	certFile = flag.String("cert", "", "")
	keyFile  = flag.String("key", "", "")
	root     = flag.String("root", "", "")
	port     = flag.Int("port", -1, "")
	realm    = flag.String("realm", "", "")
	htdigest = flag.String("htdigest", "", "")
)

func makeAuthHandler(handler http.Handler) auth.AuthenticatedHandlerFunc {
	return func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		handler.ServeHTTP(w, &r.Request)
	}
}

func main() {
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		if fmt.Sprintf("%s", f.Value) == f.DefValue {
			log.Println("Arguments invalid")
			os.Exit(0)
		}
	})

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

	secret := auth.HtdigestFileProvider(*htdigest)
	authenticator := auth.NewDigestAuthenticator(*realm, secret)

	log.Printf("Starting server on port %d\n", *port)

	hostport := fmt.Sprintf(":%d", *port)
	authHandler := makeAuthHandler(http.FileServer(http.Dir(absRoot)))
	handler := authenticator.Wrap(authHandler)
	log.Fatal(http.ListenAndServeTLS(hostport, absCert, absKey, handler))
}
