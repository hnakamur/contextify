package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hnakamur/contextify"
)

func main() {
	runTime := flag.Duration("run", 5*time.Second, "duration before run finishes without cancel")
	shutdownTime := flag.Duration("shutdown", 1*time.Second, "duration for shtudown")
	triggerTime := flag.Duration("trigger", time.Millisecond, "duration for triggering shtudown")
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

	shutdown := make(chan struct{})
	run := contextify.Contextify(func() error {
		log.Printf("run func started")
		t := time.NewTimer(*runTime)
		defer t.Stop()
		select {
		case <-t.C:
			// do nothing
		case <-shutdown:
			log.Printf("received shudown in run func")
			time.Sleep(*shutdownTime)
		}
		log.Printf("exiting run func")
		return nil
	}, func() error {
		log.Printf("cancel func started")
		close(shutdown)
		log.Printf("trigger shutdown from cancel func to run func")
		time.Sleep(*triggerTime)
		log.Printf("exiting cancel func")
		return nil
	}, nil)

	err := run(ctx)
	if err != nil {
		log.Printf("got error from run, %v", err)
	}
	log.Print("exiting main")
}
