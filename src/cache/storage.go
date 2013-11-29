package cache

import (
	"fmt"
	"net/http"
	"time"
)

const (
	STATUS_NORMAL   = 1
	STATUS_UPDATING = 2
)

type ResponseStorage struct { //响应体
	Header     http.Header
	Body       []byte
	StatusCode int
}

type Storage struct {
	InitAt             time.Time     //初始化时间
	UpdatedAt          time.Time     //最后更新时间
	Duration           int64         //(秒)缓存持续时间,如最后更新时间 - 初始化时间 > Duration 则需要更新
	ClientLastAccessAt time.Time     //客户端最后访问时间
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
	s, exists = cacheStorage[k]
	if exists {
		s.ClientAccessCount++
		s.ClientLastAccessAt = time.Now()
	}
	return
}

func (c *Cache) Set(k string, v *Storage) {
	cacheStorage[k] = v
}

func (c *Cache) Size() int {
	return len(cacheStorage)
}

/**
返回需要更新的缓存内容
*/
func (c *Cache) TimeoutEntities() (rs []*Storage) {
	var s *Storage
	for k, _ := range cacheStorage {
		s = cacheStorage[k]
		if (time.Now().Unix() - s.UpdatedAt.Unix()) > s.Duration {
			rs = append(rs, s)
		}
	}
	return
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
