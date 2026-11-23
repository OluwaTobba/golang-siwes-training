package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")
    })
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from a Go container! DB host: %s\n", os.Getenv("DB_HOST"))
    })
    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
