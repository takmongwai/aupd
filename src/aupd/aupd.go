package main

import (
  "log"
  "os"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  log.Println("Start serving on port 8000")
  http.ListenAndServe(":8000", nil)
  os.Exit(0)
}
