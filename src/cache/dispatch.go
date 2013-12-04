package cache

import (
  "client"
  "log"
  _ "net/http"
  "time"
)

var c = New()
var is_running bool = false

/*
更新超时的实体
*/

func update_timeout_entity(s *Storage) (err error) {
  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in update_timeout_entity:", re)
    }
  }()
  defer func() { s.CurrentStatus = STATUS_NORMAL }()

  s.CurrentStatus = STATUS_UPDATING
  r := s.Request
  body, status_code, header, err := client.HttpRequestNotResponse(r)
  if err != nil {
    log.Printf("update_timeout_entity error. %v", err)
    return
  }

  if status_code != 200 {
    log.Printf("update %s,status %d\n", s.Request.URL.String(), status_code)
    return
  }
  s.UpdatedAt = time.Now()
  s.Response.Body = body
  s.Response.StatusCode = status_code
  s.Response.Header = *header
  return
}

func Dispatch() {
  errc := make(chan error, MAX_CONCURRENT)
  quit := make(chan struct{})
  defer close(quit)

  for {
    time.Sleep(time.Millisecond * 500)
    c.RemoveOldEntities()
    ts := c.TimeoutEntities()
    if len(ts) > 0 {
      log.Println("begin update ", ts)
    }
    for _, s := range ts {
      go func(s *Storage) {
        select {
        case errc <- update_timeout_entity(s):
          log.Printf("update %s done", s.Request.URL.String())
        case <-quit:
          log.Printf("update %s quit", s.Request.URL.String())
        }
      }(s)
    }
    for _ = range ts {
      if err := <-errc; err != nil {
        log.Println(err)
      }
    }
    if len(ts) > 0 {
      log.Println("end update ", ts)
    }
  }
}
