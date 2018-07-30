package redis

import (
	"github.com/go-redis/redis"
	"brief_framework/logger"
	"time"
)

type Client struct {
	*redis.Client
}

const Nil = redis.Nil

func init() {

}

func NewRedisClient(addr, passwd string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	c := Client{Client: client}
	return &c, nil
}

func (c *Client) Set(key, value string) error {
	err := c.Client.Set(key, value, 0).Err()
	if err != nil {
		logger.Instance().Warn("redis set error, err = %v\n", err)
		return err
	}
	return nil
}

func (c *Client) SetExpire(key, value string, expire int64) error {
	err := c.Client.Set(key, value, time.Duration(expire)).Err()
	if err != nil {
		logger.Instance().Warn("redis set error, err = %v\n", err)
		return err
	}
	return nil
}

func (c *Client) Get(key string) (string, error) {
	ret, err := c.Client.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Instance().Warn("redis get error, err = %v\n", err)
		}
		return "", nil
	}
	return ret, nil
}

func (c *Client) Exist(key string) (bool, error) {

	val, err := c.Client.Exists(key).Result()
	if err != nil {
		logger.Instance().Warn("redis Exist error, err = %v\n", err)
		return false, err
	}
	if val > 0 {
		return true, nil
	}
	return false, nil

}

func (c *Client) HSet(key, field, value string) error {
	err := c.Client.HSet(key, field, value).Err()
	if err != nil {
		logger.Instance().Warn("redis HSet error, err = %v\n", err)
		return err
	}
	return nil
}

func (c *Client) HGet(key, field string) (string, error) {
	ret, err := c.Client.HGet(key, field).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Instance().Warn("redis hget error, err = %v\n", err)
		}
		return "", err
	}
	return ret, nil
}

func (c *Client) HGetAll(key string) (map[string]string, error) {
	ret, err := c.Client.HGetAll(key).Result()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis HGetAll error, err = %v\n", err)
		//}
		return nil, err
	}
	return ret, nil
}

func (c *Client) Keys(pattern string) (*[]string, error) {
	ret, err := c.Client.Keys(pattern).Result()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis Keys error, err = %v\n", err)
		//}
		return nil, err
	}
	return &ret, nil
}

func (c *Client) ZRangeByScoreWithScores(key string, start, stop string, offset, count int64) (*[]redis.Z, error) {
	ret, err := c.Client.ZRangeByScoreWithScores(key, redis.ZRangeBy{start, stop, offset, count}).Result()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis ZRangeWithScores error, err = %v\n", err)
		//}
		return nil, err
	}
	return &ret, nil
}

func (c *Client) ZRangeWithScores(key string, start, stop int64) (*[]redis.Z, error) {
	ret, err := c.Client.ZRangeWithScores(key, start, stop).Result()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis ZRangeWithScores error, err = %v\n", err)
		//}
		return nil, err
	}
	return &ret, nil
}

func (c *Client) PrintZRangeWithScores(key string, start, stop int64) (string) {
	return c.Client.ZRangeWithScores(key, start, stop).String()
}

func (c *Client) ZRange(key string, start, stop int64) (*[]string, error) {
	ret, err := c.Client.ZRange(key, start, stop).Result()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis ZRange error, err = %v\n", err)
		//}
		return nil, err
	}
	return &ret, nil
}

func (c *Client) ZAdd(key string, score int64, mem string) error {
	member := redis.Z{float64(score), mem}
	err := c.Client.ZAdd(key, member).Err()
	if err != nil {
		//if err != redis.Nil {
		logger.Instance().Warn("redis ZRange error, err = %v\n", err)
		//}
		return err
	}
	return nil
}