package contextify_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/hnakamur/contextify"
)

func ExampleContextify() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		s := <-c
		log.Printf("received signal, %s", s)
		cancel()
		log.Printf("cancelled context")
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, example http server\n"))
	})
	s := http.Server{Addr: ":8080"}
	run := contextify.Contextify(func() error {
		return s.ListenAndServe()
	}, func() error {
		return s.Shutdown(context.Background())
	}, nil)
	err := run(ctx)
	if err != nil {
		log.Printf("got error, %v", err)
	}
	log.Print("exiting")
}
