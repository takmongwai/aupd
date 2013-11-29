package cache

import (
  _ "errors"
  "io"
)

func Copy(dst io.Writer, src io.Reader) (all_buf []byte, written int64, err error) {

  buf := make([]byte, 32*1024)
  for {
    nr, er := src.Read(buf)
    if nr > 0 {
      nw, ew := dst.Write(buf[0:nr])
      if nw > 0 {
        written += int64(nw)
      }
      if ew != nil {
        err = ew
        break
      }
      if nr != nw {
        err = io.ErrShortWrite
        break
      }
      all_buf = append(all_buf, buf[0:nr]...)
    }
    if er == io.EOF {
      break
    }
    if er != nil {
      err = er
      break
    }
  }
  if int64(len(all_buf)) != written {
    err = io.ErrShortWrite
  }
  return all_buf, written, err
}
