package main

import (
  "fmt"
  "time"
)

type Storage struct {
  InitAt time.Time
}

type Cache map[string]*Storage

var c = make(Cache)
var s *Storage

func main() {
  fmt.Println("================================")

  for i := 0; i <= 10; i++ {
    s = &Storage{InitAt: time.Now()}
    c[fmt.Sprintf("key_%d", i)] = s
  }
  fmt.Println(c)
  
}
