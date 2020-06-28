package services

import "github.com/go-redis/redis"

// Set 集合
type Set interface {
	Exists(v interface{}) bool
	Add(v interface{}) error
	Remove(v interface{}) error
}

// RedisSet 使用 Redis 实现集合
type RedisSet struct {
	key    string
	client *redis.Client
}

// NewRedisSet 初始化基于 Redis 的集合
func NewRedisSet(k string, c *redis.Client) *RedisSet {
	return &RedisSet{
		key:    k,
		client: c,
	}
}

// Exists 在 Redis 集合中检查是否存在元素
func (s *RedisSet) Exists(v interface{}) bool {
	cmd := s.client.SIsMember(s.key, v)
	return cmd.Val()
}

// Add 向 Redis 的集合中插入元素
func (s *RedisSet) Add(v interface{}) error {
	return s.client.SAdd(s.key, v).Err()
}

// Remove 删除元素
func (s *RedisSet) Remove(v interface{}) error {
	return s.client.SRem(s.key, v).Err()
}

var _ Set = (*RedisSet)(nil)
