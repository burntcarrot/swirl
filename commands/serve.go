package commands

import (
	"fmt"
	"log"
	"net/http"
)

func requestLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("%s\t%s", r.Method, r.URL.Path)
	})
}

func Serve(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./build")))

	server := http.Server{
		Addr:    addr,
		Handler: requestLog(mux),
	}

	fmt.Println("serving on", addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
