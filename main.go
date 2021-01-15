package utils

import (
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// ListenAddrs Listen http and https
func ListenAddrs(addr, addTLS, cert, key string, handler http.Handler) {
	var g errgroup.Group
	if addTLS != "" {
		g.Go(func() error {
			return http.ListenAndServeTLS(addTLS, cert, key, handler)
		})
	}
	if addr != "" {
		g.Go(func() error { return http.ListenAndServe(addr, handler) })
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
