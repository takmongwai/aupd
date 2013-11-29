package main

import (
  "log"
  "net/http"
  "os"
  _"time"
  "cache"
)


func main() {
  go cache.Dispatch()
  http.HandleFunc("/", handler)
  log.Println("Start serving on port 8000")
  http.ListenAndServe(":8000", nil)
  os.Exit(0)
}
