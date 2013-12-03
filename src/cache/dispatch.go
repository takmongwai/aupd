package cache

import (
  "client"
  "log"
  _ "net/http"
  "time"
)

var c = New()

/*
更新超时的实体
*/

func update_timeout_entity(s *Storage) {
  defer func() {
    if re := recover(); re != nil {
      log.Println("Recovered in update_timeout_entity:", re)
    }
  }()
  start_time := time.Now()
  s.CurrentStatus = STATUS_UPDATING
  r := s.Request
  body, status_code, header, err := client.HttpRequestNotResponse(r)
  if err != nil {
    log.Printf("update_timeout_entity error. %v", err)
    return
  }
  defer func() { s.CurrentStatus = STATUS_NORMAL }()
  if status_code != 200 {
    log.Printf("update %s,status %d\n", s.Request.URL.String(), status_code)
    return
  }
  s.UpdatedAt = time.Now()
  s.Response.Body = body
  s.Response.StatusCode = status_code
  s.Response.Header = *header
  /*
  log.Printf(`update "%s",%d,[%s],%v Sec`,
    s.Request.URL.String(),
    s.ClientAccessCount,
    s.ClientLastAccessAt.Format("2006-01-02 15:04:05"),
    time.Now().Sub(start_time).Seconds(),
  )*/
}

func Dispatch() {
  for {
    time.Sleep(time.Millisecond * 1000)
    c.RemoveOldEntities()
    for _, s := range c.TimeoutEntities() {
      go update_timeout_entity(s)
    }
  }
}
