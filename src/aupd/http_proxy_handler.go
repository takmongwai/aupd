package main

import (
  _ "cache"
  "io"
  _ "io/ioutil"
  "log"
  "net"
  "net/http"
  "time"
)

var transport = http.Transport{
  Dial: func(nework, addr string) (net.Conn, error) {
    return net.DialTimeout(nework, addr, time.Duration(20)*time.Second)
  },
  ResponseHeaderTimeout: time.Duration(60) * time.Second,
  DisableCompression:    false,
  DisableKeepAlives:     false,
  MaxIdleConnsPerHost:   200,
}

var client = &http.Client{
    Transport: &transport,
  }

func headerCopy(s http.Header, d *http.Header) {
  for hk, _ := range s {
    d.Set(hk, s.Get(hk))
  }
}

func showError(w http.ResponseWriter, msg []byte) {
  w.WriteHeader(500)
  w.Write(msg)
}

func handler(w http.ResponseWriter, r *http.Request) {
  
  log.Println("----------------------------------")
  log.Println("RequestURI", r.RequestURI)
  log.Println("RemoteAddr", r.RemoteAddr)
  log.Println("----------------------------------")

  req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
  headerCopy(r.Header, &req.Header)
  defer func() { req.Close = true }()
  if err != nil {
    showError(w, []byte(err.Error()))
    return
  }
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    showError(w, []byte(err.Error()))
    return
  }
  for hk, _ := range resp.Header {
    w.Header().Set(hk, resp.Header.Get(hk))
  }
  w.WriteHeader(resp.StatusCode)
  io.Copy(w, resp.Body)

}
