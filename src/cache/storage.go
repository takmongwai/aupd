package cache

import (
  "fmt"
  "log"
  "net/http"
  "sync"
  "time"
)

var lock = sync.Mutex{}

const (
  STATUS_NORMAL               = 1
  STATUS_UPDATING             = 2
  ENTITY_UPDATE_DURATION      = 5    //Second,每个缓存需要更新的时间,依据最后更新时间和当前时间计算
  ENTITY_DURATION             = 3600 
  CLIENT_LAST_ACCESS_DURATION = 3600 //每个缓存的持续时间,依据最后访问时间和当前时间计算
)

type ResponseStorage struct { //响应体
  Header     http.Header
  Body       []byte
  StatusCode int
}

type Storage struct {
  InitAt             time.Time     //初始化时间
  UpdatedAt          time.Time     //最后更新时间
  UpdateDuration     int64         //(秒)缓存更新时间: if now() - UpdatedAt > this then remov
  Duration           int64         //(秒)缓存持续时间: now() - InitAt > this then remov
  ClientLastAccessAt time.Time     //客户端最后访问时间: 如果连续1个小时没有客户端访问,则删除
  ClientAccessCount  int           //客户端访问次数
  CurrentStatus      int           //当前状态
  Request            *http.Request //向后端发请求的结构
  Response           *ResponseStorage
}

type Cache map[string]*Storage

var cacheStorage = make(Cache)

var _instance Cache

func New() Cache {
  if _instance == nil {
    _instance = make(Cache)
  }
  return _instance
}

func (c *Cache) Get(k string) (s *Storage, exists bool) {
  lock.Lock()
  defer lock.Unlock()
  s, exists = cacheStorage[k]
  if exists {
    s.CurrentStatus = STATUS_UPDATING
    s.ClientAccessCount++
    s.ClientLastAccessAt = time.Now()
    s.CurrentStatus = STATUS_NORMAL
  }
  return
}

func (c *Cache) Exists(k string) (exists bool) {
  _, exists = cacheStorage[k]
  return
}

func (c *Cache) Set(k string, v *Storage) {
  lock.Lock()
  defer lock.Unlock()
  cacheStorage[k] = v
}

func (c *Cache) Size() int {
  return len(cacheStorage)
}

/**
返回需要更新的缓存内容
*/
func (c *Cache) TimeoutEntities() (rs []*Storage) {
  lock.Lock()
  defer lock.Unlock()
  for _, s := range cacheStorage {
    if (time.Now().Unix() - s.UpdatedAt.Unix()) > s.UpdateDuration {
      rs = append(rs, s)
    }
  }
  return
}

/**
删除不再需要的缓存
1: 最后一次客户端访问时间距离当前时间超过1个小时的
*/
func (c *Cache) RemoveOldEntities() {
  lock.Lock()
  defer lock.Unlock()
  for k, s := range cacheStorage {
    if time.Now().Unix()-s.ClientLastAccessAt.Unix() > CLIENT_LAST_ACCESS_DURATION {
      log.Printf("RemoveOldEntities: %s\n", k)
      delete(cacheStorage, k)
    }
  }
}

func (c *Cache) Remove(k string) {
  delete(cacheStorage, k)
}
func (s *Storage) Info() string {
  rs := fmt.Sprintf(
    `
    InitAt: %s
    UpdatedAt: %s
    Duration: %d
    ClientLastAccessAt: %s
    ClientAccessCount: %d
    CurrentStatus: %d
    `,
    s.InitAt,
    s.UpdatedAt,
    s.Duration,
    s.ClientLastAccessAt,
    s.ClientAccessCount,
    s.CurrentStatus,
  )
  return rs
}
