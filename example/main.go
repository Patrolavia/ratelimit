package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	var server, client bool
	flag.BoolVar(&server, "server", false, "Run a web server which limit transfer rate to 10k/s")
	flag.BoolVar(&client, "client", false, "Demo client transfer rate limitation, cannot use with -server")
	flag.Parse()

	if server == client {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if server {
		log.Print("Starting server as :8000 ...")
		log.Print("Access http://127.0.0.1:8000/ to test global rate limit.")
		log.Print("http://127.0.0.1:8000/10k to test per-thread rate limit.")
		if err := serverExample(); err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	} else {
		mockserver()
		time.Sleep(time.Second)
		log.Print("Server listening to :8000, ready to download.")
		beginTime := time.Now()
		if err := clientExample(); err != nil {
			log.Fatalf("Cannot download data: %s", err)
		}
		endTime := time.Now()
		dur := float64(endTime.Sub(beginTime) / time.Second)
		log.Printf("Downloaded 100k file in %f seconds, %fkb/s", dur, 100.0/dur)
	}
}
