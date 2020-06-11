package services

import (
	"bytes"
	"fmt"
	"sync"
	"time"
	"translator-api/app"
	"translator-api/hash"

	"github.com/go-redis/redis"
	"github.com/golang/groupcache/lru"
	jsoniter "github.com/json-iterator/go"
	"github.com/zgs225/youdao"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// YoudaoService 查询有道数据
type YoudaoService interface {
	Query(string) (*youdao.Result, error)
}

type basicYoudaoService struct {
	pool sync.Pool
}

// YoudaoServiceMiddleware 有道中间件
type YoudaoServiceMiddleware func(YoudaoService) YoudaoService

// NewBasicYoudaoService 获取有道查询服务
func NewBasicYoudaoService(appID, appSecret string) YoudaoService {
	return &basicYoudaoService{
		pool: sync.Pool{
			New: func() interface{} {
				return &youdao.Client{
					AppID:     appID,
					AppSecret: appSecret,
				}
			},
		},
	}
}

// Query 使用有道 API 查询
func (s *basicYoudaoService) Query(q string) (*youdao.Result, error) {
	client := s.pool.Get().(*youdao.Client)
	return client.Query(q)
}

// RedisCachedYoudaoService 保存在 Redis 中的有道查询
type RedisCachedYoudaoService struct {
	next YoudaoService
	rds  *redis.Client
	core *app.Application
}

// NewRedisCachedYoudaoServiceMiddleware 生成使用Redis缓存的服务中间件
func NewRedisCachedYoudaoServiceMiddleware(r *redis.Client, core *app.Application) YoudaoServiceMiddleware {
	return func(next YoudaoService) YoudaoService {
		return &RedisCachedYoudaoService{
			next: next,
			rds:  r,
			core: core,
		}
	}
}

// Query 在 Redis 中查询
func (s *RedisCachedYoudaoService) Query(q string) (*youdao.Result, error) {
	key := fmt.Sprintf("yd:v1:%s", hash.SHA256(q))
	str, err := s.rds.Get(key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(str) == 0 {
		s.core.Logger.Debug("Redis cache miss: ", q, " by key: ", key)
		result, err := s.next.Query(q)
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(result); err != nil {
			s.core.Logger.WithError(err).Error("JSON serialize error.")
			return result, nil
		}
		if err := s.rds.Set(key, buf.String(), 365*24*time.Hour).Err(); err != nil {
			s.core.Logger.WithError(err).Error("Redis store result error.")
		}
		return result, nil
	}

	s.core.Logger.Debug("Redis cache hit: ", q, " by key: ", key)
	result := youdao.Result{}
	err = json.Unmarshal([]byte(str), &result)
	return &result, err
}

// MemoryCachedYoudaoService 内存缓存
type MemoryCachedYoudaoService struct {
	cache *lru.Cache
	core  *app.Application
	next  YoudaoService
}

// NewMemoryCachedYoudaoService 内存缓存中间件
func NewMemoryCachedYoudaoService(cache *lru.Cache, core *app.Application) YoudaoServiceMiddleware {
	return func(next YoudaoService) YoudaoService {
		return &MemoryCachedYoudaoService{
			cache: cache,
			core:  core,
			next:  next,
		}
	}
}

// Query 在内存中查询
func (s *MemoryCachedYoudaoService) Query(q string) (*youdao.Result, error) {
	key := "yd:lru:" + hash.SHA256(q)
	v, ok := s.cache.Get(lru.Key(key))
	if ok {
		s.core.Logger.Debug("LRU cache hit: ", q, " by key: ", key)
		return v.(*youdao.Result), nil
	}
	s.core.Logger.Debug("LRU cache miss: ", q, " by key: ", key)
	result, err := s.next.Query(q)
	if err != nil {
		return nil, err
	}
	s.cache.Add(key, result)
	return result, nil
}
