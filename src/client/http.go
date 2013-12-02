package client

import (
	"io/ioutil"
	_ "log"
	"net"
	"net/http"
	"sync"
	"time"
	"util"
)

const (
	CONNECTION_TIME_OUT     = 5  //连接超时
	RESPONSE_TIME_OUT       = 20 //响应超时
	MAX_IDLE_CONNS_PRE_HOST = 200
	DISABLE_COMPRESSION     = false
	DISABLE_KEEP_ALIVES     = false
	MAX_CACHE_ENTITY        = 1024 * 512 //byte
)

var lock = sync.Mutex{}

type HttpResponse struct { //响应体
	Header     http.Header
	Body       []byte
	StatusCode int
}

var transport = http.Transport{
	Dial: func(nework, addr string) (net.Conn, error) {
		return net.DialTimeout(nework, addr, time.Duration(CONNECTION_TIME_OUT)*time.Second)
	},
	ResponseHeaderTimeout: time.Duration(RESPONSE_TIME_OUT) * time.Second,
	DisableCompression:    DISABLE_COMPRESSION,
	DisableKeepAlives:     DISABLE_KEEP_ALIVES,
	MaxIdleConnsPerHost:   MAX_IDLE_CONNS_PRE_HOST,
}

var client = &http.Client{
	Transport: &transport,
}

func headerCopy(s http.Header, d *http.Header) {
	lock.Lock()
	for hk, _ := range s {
		d.Set(hk, s.Get(hk))
	}
	lock.Unlock()
}

func showError(w http.ResponseWriter, msg []byte, outbuf []byte, written *int64) {
	outbuf = msg
	*written = int64(len(msg))
	w.WriteHeader(500)
	w.Write(msg)
}

func HttpRequestNotResponse(r *http.Request) (body []byte, resp_status_code int, resp_header *http.Header, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		panic(err)
	}
	headerCopy(r.Header, &req.Header)
	defer func() { req.Close = true }()
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	resp_status_code = resp.StatusCode

	resp_header = &http.Header{}
	headerCopy(resp.Header, resp_header)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return
}

func HttpRequest(w http.ResponseWriter, r *http.Request) (body []byte, resp_status_code int, written int64, err error) {
	var req *http.Request

	req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
	headerCopy(r.Header, &req.Header)
	defer func() { req.Close = true }()

	if err != nil {
		showError(w, []byte(err.Error()), body, &written)
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		showError(w, []byte(err.Error()), body, &written)
		return
	}

	for hk, _ := range resp.Header {
		w.Header().Set(hk, resp.Header.Get(hk))
	}

	resp_status_code = resp.StatusCode
	w.WriteHeader(resp_status_code)

	body, written, err = util.Copy(w, resp.Body)
	if err != nil {
		showError(w, []byte(err.Error()), body, &written)
		return
	}
	return
}
