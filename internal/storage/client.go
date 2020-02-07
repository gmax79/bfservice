package storage

import "github.com/go-redis/redis/v7"

// Provider - load/save settings in extenal storage
type Provider interface {
	CreateSet(id string) (SetProvider, error)
	Close() error
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
	status := p.rc.Ping()
	return &p, status.Err()
}

func (p *redisProvider) CreateSet(id string) (SetProvider, error) {
	return createRedisSetProvider(p.rc, id)
}

func (p *redisProvider) Close() error {
	return p.rc.Close()
}
