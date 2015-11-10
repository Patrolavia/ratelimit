package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Patrolavia/ratelimit"
)

func serverExample() (err error) {
	bucket := ratelimit.NewFromRate(10*1024, 10*1024, 0)
	content := strings.Repeat(".", 100*1024)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// all server threads share 10k bandwidth by using same bucket
		wrappedWriter := ratelimit.NewWriter(w, bucket)
		fmt.Fprint(wrappedWriter, content)
	})
	http.HandleFunc("/10k", func(w http.ResponseWriter, r *http.Request) {
		// each thread has 10k bandwidth
		bucket := ratelimit.NewFromRate(10*1024, 10*1024, 0)
		wrappedWriter := ratelimit.NewWriter(w, bucket)
		fmt.Fprint(wrappedWriter, content)
	})

	return http.ListenAndServe(":8000", nil)
}
