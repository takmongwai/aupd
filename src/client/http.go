package client

import (
  "fmt"
  "io/ioutil"
  "log"
  "net"
  "net/http"
_  "strings"
  "sync"
  "time"
  "util"
)

const (
  CONNECTION_TIME_OUT     = 15 //连接超时,秒
  RESPONSE_TIME_OUT       = 90 //响应超时,秒
  MAX_IDLE_CONNS_PRE_HOST = 6
  DISABLE_COMPRESSION     = false
  DISABLE_KEEP_ALIVES     = true
  MAX_CACHE_ENTITY        = 1024 * 512 //byte
)

var lock = sync.Mutex{}

type HttpResponse struct { //响应体
  Header     http.Header
  Body       []byte
  StatusCode int
}

func timeoutDialer(conn_timeout int, rw_timeout int) func(net, addr string) (c net.Conn, err error) {
  return func(netw, addr string) (net.Conn, error) {
    conn, err := net.DialTimeout(netw, addr, time.Duration(conn_timeout)*time.Second)
    if err != nil {
      log.Printf("Failed to connect to [%s]. Timed out after %d seconds\n", addr, rw_timeout)
      return nil, err
    }
    conn.SetDeadline(time.Now().Add(time.Duration(rw_timeout) * time.Second))
    return conn, nil
  }
}

var transport = http.Transport{
  Dial: timeoutDialer(CONNECTION_TIME_OUT, RESPONSE_TIME_OUT),
  ResponseHeaderTimeout: time.Duration(RESPONSE_TIME_OUT) * time.Second,
  DisableCompression:    DISABLE_COMPRESSION,
  DisableKeepAlives:     DISABLE_KEEP_ALIVES,
  MaxIdleConnsPerHost:   MAX_IDLE_CONNS_PRE_HOST,
}

var client = &http.Client{
  Transport: &transport,
}

func cleanHeader(h *http.Header) {
  if len(h.Get("Accept-Encoding")) > 0 {
    h.Del("Accept-Encoding")
  }
}

func headerCopy(s http.Header, d *http.Header) {
  lock.Lock()
  defer lock.Unlock()
  for hk, _ := range s {
    d.Set(hk, s.Get(hk))
  }
}

func showError(w http.ResponseWriter, msg []byte, outbuf []byte, written *int64) {
  outbuf = msg
  *written = int64(len(msg))
  for hk, _ := range w.Header() {
    w.Header().Del(hk)
  }
  w.WriteHeader(500)
  w.Write(msg)
}

func HttpRequestNotResponse(r *http.Request) (body []byte, resp_status_code int, resp_header *http.Header, err error) {
  var (
    req  *http.Request
    resp *http.Response
  )
  req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
  if err != nil {
    panic(err)
  }
  headerCopy(r.Header, &req.Header)
  cleanHeader(&req.Header)
  defer func() { req.Close = true }()
  resp, err = client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  resp_status_code = resp.StatusCode

  resp_header = &http.Header{}
  headerCopy(resp.Header, resp_header)
  body, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  return
}

func HttpRequest(w http.ResponseWriter, r *http.Request) (body []byte, resp_status_code int, written int64, err error) {
  var req *http.Request

  req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
  headerCopy(r.Header, &req.Header)
  cleanHeader(&req.Header)
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

  resp_status_code = resp.StatusCode
  w.WriteHeader(resp_status_code)

  body, written, err = util.Copy(w, resp.Body)

  if err != nil {
    showError(w, []byte(err.Error()), body, &written)
    return
  }
  return
}

func FullQueryString(r *http.Request) (rs string) {
  var rawQuery []string
  r.ParseForm()
  for k, _ := range r.Form {
    rawQuery = append(rawQuery, fmt.Sprintf("%s=%s", k, r.Form.Get(k)))
  }
  rs = r.RequestURI
  if len(rawQuery) > 0 {
    rs = fmt.Sprintf("%s", rs)
  }
  return
}
