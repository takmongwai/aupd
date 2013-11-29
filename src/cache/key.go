package cache

import (
  "crypto/md5"
  "fmt"
  "io"
  _ "log"
  "net/http"
  "sort"
  "strings"
)

/*
获取get,post的请求参数,返回一个string slice,并根据参数名进行排序
*/

func parseQuerys(r *http.Request) (rawQuery []string) {
  r.ParseForm()
  for k, _ := range r.Form {
    rawQuery = append(rawQuery, fmt.Sprintf("%s=%s", k, r.Form.Get(k)))
  }
  sort.Strings(rawQuery)
  return
}

/*
根据请求生成缓存 key
(URL + 排序过的参数(含post)).md5
*/
func GenKey(r *http.Request) string {
  sorted_params := strings.Join(parseQuerys(r),"&")
  return str_md5(fmt.Sprintf("%s?%s",r.RequestURI,sorted_params))
}

func str_md5(s string) string {
  h := md5.New()
  io.WriteString(h, s)
  return fmt.Sprintf("%x", h.Sum(nil))
}
