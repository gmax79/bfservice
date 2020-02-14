package storage

import "github.com/go-redis/redis/v7"

type redisSetProvider struct {
	rc *redis.Client
	id string
}

func createRedisSetProvider(rc *redis.Client, id string) (*redisSetProvider, error) {
	var p redisSetProvider
	p.rc = rc
	p.id = id
	return &p, nil
}

func (p *redisSetProvider) Add(item string) (bool, error) {
	result := p.rc.SAdd(p.id, item)
	v, err := result.Val(), result.Err()
	flag := v == 0
	return flag, err
}

func (p *redisSetProvider) Delete(item string) (bool, error) {
	result := p.rc.SRem(p.id, item)
	v, err := result.Val(), result.Err()
	flag := v == 0
	return flag, err
}

func (p *redisSetProvider) Iterator() (StringIterator, error) {
	cmd := p.rc.SScan(p.id, 0, "", 0)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	iterator := func() (string, bool) {
		val := cmd.String()
		return val, cmd.Err() != nil
	}
	return iterator, nil
}
