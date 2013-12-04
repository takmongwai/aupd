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
  //log.Printf("begin update %s,%s", s.Request.URL.String(), s.InitAt)
  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in update_timeout_entity:", re, " at ", client.FullQueryString(s.Request))
    }
    s.CurrentStatus = STATUS_NORMAL
  }()

  s.CurrentStatus = STATUS_UPDATING
  body, status_code, header, err := client.HttpRequestNotResponse(s.Request)
  if err != nil {
    log.Printf("update_timeout_entity error. %v at %s", err, client.FullQueryString(s.Request))
    return
  }

  if status_code != 200 {
    log.Printf("update %s,status %d\n", client.FullQueryString(s.Request), status_code)
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
  ts := make([]*Storage, MAX_CONCURRENT)
  for {
    time.Sleep(time.Millisecond * 300)
    seq := time.Now().UnixNano()
    c.RemoveOldEntities()
    ts = c.TimeoutEntities()

    if len(ts) > 0 {
      log.Printf("begin %d,size: %d", seq, len(ts))
      log.Println("total cached: %d", c.Size())
    }

    for i := 0; i < len(ts); i++ {
      go func(fs *Storage) {
        select {
        case errc <- update_timeout_entity(fs):
          log.Printf("update %s done", client.FullQueryString(fs.Request))
        case <-quit:
          log.Printf("update %s quit", client.FullQueryString(fs.Request))
        }
      }(ts[i])
    }

    for i := 0; i < len(ts); i++ {
      if err := <-errc; err != nil {
        log.Println(err)
      }
    }
    
    if len(ts) > 0 {
      log.Printf("end %d,size: %d", seq, len(ts))
    }

  }
}
