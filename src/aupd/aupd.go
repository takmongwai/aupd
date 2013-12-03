package main

import (
  "cache"
  "log"
  "net/http"
  "os"
  _ "time"
)

func main() {
  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in main:", re)
    }
  }()
  go cache.Dispatch()
  http.HandleFunc("/", handler)
  log.Println("Start serving on port 8000")
  http.ListenAndServe(":8000", nil)
  os.Exit(0)
}
