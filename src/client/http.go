package client

import (
  "ioutil"
  "net"
  "net/http"
  "sync"
  "time"
  "util"
)

const (
  CONNECTION_TIME_OUT     = 15
  RESPONSE_TIME_OUT       = 60
  MAX_IDLE_CONNS_PRE_HOST = 200
  DISABLE_COMPRESSION     = false
  DISABLE_KEEP_ALIVES     = false
  MAX_CACHE_ENTITY        = 1024 * 512 //byte
)

var lock = sync.Mutex{}

var transport = http.Transport{
  Dial: func(nework, addr string) (net.Conn, error) {
    return net.DialTimeout(nework, addr, time.Duration(CONNECTION_TIME_OUT)*time.Second)
  },
  ResponseHeaderTimeout: time.Duration(RESPONSE_TIME_OUT) * time.Second,
  DisableCompression:    DISABLE_COMPRESSION,
  DisableKeepAlives:     DISABLE_KEEP_ALIVES,
  MaxIdleConnsPerHost:   MAX_IDLE_CONNS_PRE_HOST,
}

var client = &http.Client{
  Transport: &transport,
}

func headerCopy(s http.Header, d *http.Header) {
  lock.Lock()
  for hk, _ := range s {
    d.Set(hk, s.Get(hk))
  }
  lock.Unlock()
}

func showError(w http.ResponseWriter, msg []byte, outbuf []byte, written *int64) {
  outbuf = msg
  *written = int64(len(msg))
  w.WriteHeader(500)
  w.Write(msg)
}

func HttpRequestByte(r *http.Request) {
  var (
    req *http.Request
    err error
  )
  req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
  headerCopy(r.Header, &req.Header)
  defer func() { req.Close = true }()
  body, err := ioutil.ReadAll(resp.Body)
}

func HttpRequest(w http.ResponseWriter, r *http.Request) (body []byte, written int64, err error) {
  var req *http.Request

  req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
  headerCopy(r.Header, &req.Header)
  defer func() { req.Close = true }()

  if err != nil {
    showError(w, []byte(err.Error()), body, &written)
    return
  }
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    showError(w, []byte(err.Error()), body, &written)
    return
  }

  for hk, _ := range resp.Header {
    w.Header().Set(hk, resp.Header.Get(hk))
  }

  w.WriteHeader(resp.StatusCode)

  body, written, err = util.Copy(w, resp.Body)
  if err != nil {
    showError(w, []byte(err.Error()), body, &written)
    return
  }
  return
}
