package main

import (
	"cache"
	"client"
	"log"
	"net/http"
	_ "strings"
	"time"
)

var Cache = cache.New()

func showError(w http.ResponseWriter, msg []byte) {
	w.WriteHeader(500)
	w.Write(msg)
}

func handler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if re := recover(); re != nil {
			log.Println("Recovered in handler:", re)
			w.WriteHeader(500)
			w.Write([]byte("BackenServer Error"))
		}
	}()

	var (
		cache_key        = cache.GenKey(r)
		cache_storage    *cache.Storage
		cache_exists     bool
		resp_body        []byte
		err              error
		resp_status_code int
	)

	log.Printf("access:%s,key: %s", r.URL.String(), cache_key)
	if r.Header.Get("ACS_RELOAD") == "true" {
		log.Printf("RELOAD %s", r.URL.String())
		Cache.Remove(cache_key)
	}

	cache_storage, cache_exists = Cache.Get(cache_key)
	if cache_exists {
		for hk, _ := range cache_storage.Response.Header {
			w.Header().Set(hk, cache_storage.Response.Header.Get(hk))
		}
		w.Write(cache_storage.Response.Body)
		return
	}

	resp_body, resp_status_code, _, err = client.HttpRequest(w, r)

	if err != nil {
		showError(w, []byte(err.Error()))
		return
	}

	if resp_status_code != 200 {
		return
	}

	cache_storage = &cache.Storage{
		InitAt:             time.Now(),
		UpdatedAt:          time.Now(),
		UpdateDuration:     cache.ENTITY_UPDATE_DURATION,
		Duration:           cache.ENTITY_DURATION,
		ClientLastAccessAt: time.Now(),
		ClientAccessCount:  1,
		CurrentStatus:      cache.STATUS_NORMAL,
		Request:            r,
		Response: &cache.ResponseStorage{
			Header:     w.Header(),
			Body:       resp_body,
			StatusCode: resp_status_code,
		},
	}
	Cache.Set(cache_key, cache_storage)
}
