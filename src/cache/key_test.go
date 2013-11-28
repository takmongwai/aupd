package cache

import (
  "testing"
  "fmt"
)

func TestGenKey(t *testing.T) {
  fmt.Println(GenKey(nil))
  //fmt.Println(str_md5("hello"))
}
