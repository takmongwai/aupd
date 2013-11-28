package cache

import (
  "crypto/md5"
  "fmt"
  "io"
  "log"
  "net/http"
  "strings"
)

func parseQuerys(r *http.Request) (rawQuery []string) {
  r.ParseForm()
  for k, _ := range r.Form {
    rawQuery = append(rawQuery, fmt.Sprintf("%s=%s", k, r.Form.Get(k)))
  }
  return
}

/*
根据请求生成缓存 key
(URL + 排序过的参数(含post)).md5
*/
func GenKey(r *http.Request) string {
  log.Println(parseQuerys(r))
  return str_md5("hello")
}

func str_md5(s string) string {
  h := md5.New()
  io.WriteString(h, s)
  return fmt.Sprintf("%x", h.Sum(nil))
}
