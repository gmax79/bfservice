package storage

import "github.com/go-redis/redis/v7"

type redisSetProvider struct {
	rc *redis.Client
	id string
}

func createRedisSetProvider(rc *redis.Client, id string) *redisSetProvider {
	var p redisSetProvider
	p.rc = rc
	p.id = id
	return &p
}

func (p *redisSetProvider) Add(item string) (bool, error) {
	result := p.rc.SAdd(p.id, item)
	v, err := result.Val(), result.Err()
	flag := v != 0
	return flag, err
}

func (p *redisSetProvider) Delete(item string) (bool, error) {
	result := p.rc.SRem(p.id, item)
	v, err := result.Val(), result.Err()
	flag := v != 0
	return flag, err
}

func (p *redisSetProvider) Load(f func(v string) error) error {
	cmd := p.rc.SScan(p.id, 0, "", 0)
	err := cmd.Err()
	if err != nil {
		return err
	}
	iter := cmd.Iterator()
	for iter.Next() {
		if err := f(iter.Val()); err != nil {
			return err
		}
	}
	return nil
}
