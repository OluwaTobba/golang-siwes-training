package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    host, _ := os.Hostname()

    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") })
    http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ready") })
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "go-api pod=%s log_level=%s db_host=%s\n", host, os.Getenv("LOG_LEVEL"), os.Getenv("DB_HOST"))
    })

    log.Printf("listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}