package storage

import "github.com/go-redis/redis/v7"

type redisListProvider struct {
	rc *redis.Client
}

func createRedisListProvider(rc *redis.Client, id string) (*redisListProvider, error) {
	var p redisListProvider
	p.rc = rc
	return &p, nil
}

func (p redisListProvider) Add(item string) (bool, error) {
	p.rc.LInsert()
}

func (p redisListProvider) Delete(item string) (bool, error) {

}

func (p redisListProvider) Count() (int, error) {

}

func (p redisListProvider) Get(index int) (string, error) {

}
