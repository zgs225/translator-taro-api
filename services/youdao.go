package services

import (
	"sync"

	"github.com/zgs225/youdao"
)

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

func (s *basicYoudaoService) Query(q string) (*youdao.Result, error) {
	client := s.pool.Get().(*youdao.Client)
	return client.Query(q)
}
