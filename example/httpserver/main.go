package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/hnakamur/contextify"
)

func main() {
	address := flag.String("address", ":8080", "http server listen address")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, example http server\n"))
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		s := <-c
		log.Printf("received signal, %s", s)
		cancel()
		log.Printf("cancelled context")
	}()

	s := http.Server{Addr: *address}
	run := contextify.Contextify(func() error {
		defer log.Printf("exiting run function")
		return s.ListenAndServe()
	}, func() error {
		defer log.Printf("exiting cancel function")
		return s.Shutdown(context.Background())
	}, nil)
	err := run(ctx)
	if err != nil {
		log.Printf("got error, %v", err)
	}
	log.Print("exiting")
}
