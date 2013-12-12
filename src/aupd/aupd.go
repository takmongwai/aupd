package main

import (
  "cache"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strconv"
  "time"
)

var g_port int = 8000
var g_host string = "0.0.0.0"
var g_pidfile string
var g_logfile string

func file_exists(name string) bool {
  if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
      return false
    }
  }
  return true
}

func main() {

  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in main:", re)
    }
  }()

  flag.Usage = func() {
    fmt.Fprintf(os.Stderr,
      "Usage: %s \n",
      os.Args[0])
    flag.PrintDefaults()
    os.Exit(2)
  }

  flag.IntVar(&g_port, "p", 8000, "listen port")
  flag.StringVar(&g_host, "h", "0.0.0.0", "listen host")
  flag.StringVar(&g_pidfile, "pid", "", "pid file")
  flag.StringVar(&g_logfile, "log", "", "log file")
  flag.Parse()

  if g_pidfile != "" && filepath.IsAbs(g_pidfile) {
    if out, err := os.OpenFile(g_pidfile, os.O_WRONLY|os.O_CREATE, os.ModeAppend|0666); err == nil {
      out.WriteString(strconv.Itoa(os.Getpid()))
    }
  }

  if g_logfile != "" && filepath.IsAbs(g_logfile) {
    if out, err := os.OpenFile(g_logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0666); err == nil {
      log.SetOutput(out)
    }
  }

  go cache.Dispatch()

  // http.HandleFunc("/", handler)
  //   log.Printf("Start serving on %s:%d", g_host, g_port)
  //   if err := http.ListenAndServe(fmt.Sprintf("%s:%d", g_host, g_port), nil); err != nil {
  //     log.Println("ListenAndServe: ", err)
  //     os.Exit(2)
  //   }

  log.Printf("Start serving on %s:%d", g_host, g_port)
  s := &http.Server{
    Addr:           fmt.Sprintf("%s:%d", g_host, g_port),
    Handler:        http.HandlerFunc(handler),
    ReadTimeout:    30 * time.Second,
    WriteTimeout:   30 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  if err := s.ListenAndServe(); err != nil {
    log.Println("ListenAndServe: ", err)
    os.Exit(2)
  }

  os.Exit(0)

}
