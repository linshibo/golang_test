package main

import (
    "net/http"
)
func main() {
    // Simple static webserver:
    http.ListenAndServe(":8080", http.FileServer(http.Dir("./")))
}
