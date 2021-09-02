package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var port int
	var responseJSON string
	flag.IntVar(&port, "p", 8080, "Specify port")
	flag.StringVar(&responseJSON, "r", `{"data":"ok"}`, "Specify response JSON body")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("reading body %+v", err)
		}
		defer r.Body.Close()

		log.Printf("%v\t%v\t%v", r.Method, r.URL.Path, string(body))
		io.WriteString(w, responseJSON)
	})

	api := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	go func() {
		log.Printf("server started on %v", port)
		if err := api.ListenAndServe(); err != nil {
			log.Fatalf("%+v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		api.Close()
		log.Fatalf("could not stop server gracefully: %v", err)
	}
}
