package main

import (
  _ "bytes"
  "cache"
  "client"
  _ "io"
  _ "io/ioutil"
  "log"
  _ "net"
  "net/http"
  _ "sync"
  "time"
)

const (
  ENTITY_DURATION = 10 //Second
)

var Cache = cache.New()

func showError(w http.ResponseWriter, msg []byte) {
  w.WriteHeader(500)
  w.Write(msg)
}

func handler(w http.ResponseWriter, r *http.Request) {
  backend_server(w, r)
}

func backend_server(w http.ResponseWriter, r *http.Request) {
  var (
    cache_key     = cache.GenKey(r)
    cache_storage *cache.Storage
    cache_exists  bool
    resp_body     []byte
    err           error
  )

  cache_storage, cache_exists = Cache.Get(cache_key)
  if cache_exists {
    for hk, _ := range cache_storage.Response.Header {
      w.Header().Set(hk, cache_storage.Response.Header.Get(hk))
    }
    w.Write(cache_storage.Response.Body)
    return
  }

  resp_body, _, err = client.HttpRequest(w, r)
  if err != nil {
    showError(w, []byte(err.Error()))
    return
  }

  cache_storage = &cache.Storage{
    InitAt:             time.Now(),
    UpdatedAt:          time.Now(),
    Duration:           ENTITY_DURATION,
    ClientLastAccessAt: time.Now(),
    ClientAccessCount:  1,
    CurrentStatus:      cache.STATUS_NORMAL,
    Request:            r,
    Response: &cache.ResponseStorage{
      Header: w.Header(),
      Body:   resp_body,
    },
  }
  Cache.Set(cache_key, cache_storage)
}
