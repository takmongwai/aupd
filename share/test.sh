#!/bin/bash
for i in $(
echo "http://www.ibm.com"
echo "http://www.php.net"
echo "http://www.sohu.com/"
echo "http://www.sina.com.cn/"
echo "http://www.google.com.hk/"
echo "http://www.google.com/search?q=hello"
echo "http://www.163.com/"
echo "http://www.ngacn.cc/"
echo "http://www.oschina.net/"
echo "http://www.apache.org/"
); do
  curl --proxy http://127.0.0.1:8000 "${i}" -o /dev/null &
  #sleep 1
done
