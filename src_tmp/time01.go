package main

import (
  "fmt"
  "github.com/weidewang/go-strftime"
  "strconv"
  "time"
)

func main() {
  fmt.Println("======================")
  t1 := time.Now()
  time.Sleep(time.Second * 5)
  t2 := time.Now()
  fmt.Println(strconv.Atoi(strftime.Strftime(&t1, "%s")))
  fmt.Println(strconv.Atoi(strftime.Strftime(&t2, "%s")))
}
