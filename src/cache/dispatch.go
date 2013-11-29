package cache

import (
  "log"
  "net/http"
  "time"
  _"client"
)

var c = New()

//更新超时的实体
func update_timeout_entity(s *Storage) {
  log.Println(s.Request.URL.String())
  r := s.Request
  log.Println(w,r)
  log.Println(client.HttpRequest(w,r))
}

func Dispatch() {
  for {
    time.Sleep(time.Millisecond * 500)
    for _, s := range c.TimeoutEntities() {
      go update_timeout_entity(s)
    }
  }
}
