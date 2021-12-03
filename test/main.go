package main

import (
    "log"
    "net/http"
)

func main() { 
    // Simple static webserver:
    log.Fatal(http.ListenAndServe(":19999", http.FileServer(http.Dir("./"))))
}
