package storage

import "github.com/go-redis/redis/v7"

// Provider - load/save settings in extenal storage
type Provider interface {
	CreateList(id string) (ListProvider, error)
}

type redisProvider struct {
	rc *redis.Client
}

// ConnectRedis - connect to stograge
func ConnectRedis(host, password string, db int) (Provider, error) {
	var p redisProvider
	options := redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	}
	p.rc = redis.NewClient(&options)
	return &p, nil
}

func (p *redisProvider) CreateList(id string) (ListProvider, error) {
	return createRedisListProvider(p.rc, id)
}
