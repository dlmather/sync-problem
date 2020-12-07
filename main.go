package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Record struct {
	UUID  string `json:"uuid"`
	Index int    `json:"index"`
	State string `json:"state"`
}

type RecordStore map[string]Record

func main() {
	var hardmode bool
	rand.Seed(time.Now().Unix())
	rs := make(RecordStore)
	s := http.NewServeMux()
	s.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	})
	s.HandleFunc("/hardmode", func(w http.ResponseWriter, r *http.Request) {
		hardmode = !hardmode
		fmt.Fprintf(w, "hardmode: %v", hardmode)
	})
	s.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var (
			i    int
			skip interface{}
		)
		for dec.More() {
			var rec Record
			_ = dec.Decode(&skip)
			if err := dec.Decode(&rec); err != nil {
				http.Error(w, fmt.Sprintf("error parsing JSON input %s", err.Error()), http.StatusBadRequest)
				return
			}
			if rand.Float32() > 0.5 {
				rec.State = "corrupt"
			}
			rs[rec.UUID] = rec
			i++
		}
		fmt.Fprintf(w, "seeded %d records, but some have been corrupted!", i)
	})
	s.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		var rec Record
		if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
			http.Error(w, fmt.Sprintf("error parsing JSON input %s", err.Error()), http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			record, ok := rs[rec.UUID]
			if !ok {
				http.Error(w, "no uuid provided or no matching record", http.StatusBadRequest)
				return
			}
			if err := json.NewEncoder(w).Encode(record); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			if hardmode && (rand.Float32() > 0.7) {
				select {
				case <-time.NewTicker(30 * time.Second).C:
				case <-r.Context().Done():
					return
				}
			}
			rs[rec.UUID] = rec
			fmt.Fprintf(w, "recorded %v\n", rec)
			return
		default:
			http.Error(w, fmt.Sprintf("only supports GET and POST, not %s", r.Method), http.StatusBadRequest)
			return
		}
	})
	log.Println("starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
