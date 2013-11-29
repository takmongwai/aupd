package cache

import (
  "client"
  "log"
  _ "net/http"
  "time"
)

var c = New()

//更新超时的实体
func update_timeout_entity(s *Storage) {
  log.Printf("update begin,%s\n", s.Request.URL.String())
  s.CurrentStatus = STATUS_UPDATING
  r := s.Request
  body, header, err := client.HttpDoByte(r)
  if err != nil {
    log.Printf("update_timeout_entity error. %v", err)
    return
  }
  defer func() { s.CurrentStatus = STATUS_NORMAL }()
  s.UpdatedAt = time.Now()
  s.Response.Body = body
  s.Response.Header = *header
  log.Printf("update end,%s\n", s.Request.URL.String())
}

func Dispatch() {
  for {
    time.Sleep(time.Millisecond * 1000)
    log.Printf("total entities: %d", c.Size())
    for _, s := range c.TimeoutEntities() {
      go update_timeout_entity(s)
    }
  }
}
