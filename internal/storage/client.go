package storage

import "github.com/go-redis/redis/v7"

type provider struct {
	rc *redis.Client
}

func connect(host, password string, db int) (*provider, error) {
	var p provider
	options := redis.Options{
		Addr:     host,
		Password: password,
		DB:       db
	}
	p.client = redis.NewClient(&options)
	

}
