package config

import (
	"github.com/go-admin-team/go-admin-core/storage"
	"github.com/go-admin-team/go-admin-core/storage/queue"
	"github.com/go-redis/redis/v7"
	"github.com/robinjoseph08/redisqueue/v2"
	"time"
)

type Queue struct {
	Redis  *QueueRedis
	Memory *QueueMemory
}

type QueueRedis struct {
	RedisConnectOptions
	Producer *redisqueue.ProducerOptions
	Consumer *redisqueue.ConsumerOptions
}

type QueueMemory struct {
	PoolSize uint
}

var QueueConfig = new(Queue)

// Setup 启用顺序 redis > 其他 > memory
func (e Queue) Setup() (storage.AdapterQueue, error) {
	if e.Redis != nil {
		e.Redis.Consumer.ReclaimInterval = e.Redis.Consumer.ReclaimInterval * time.Second
		e.Redis.Consumer.BlockingTimeout = e.Redis.Consumer.BlockingTimeout * time.Second
		e.Redis.Consumer.VisibilityTimeout = e.Redis.Consumer.VisibilityTimeout * time.Second
		client := GetRedisClient()
		if client == nil {
			options, err := e.Redis.RedisConnectOptions.GetRedisOptions()
			if err != nil {
				return nil, err
			}
			client = redis.NewClient(options)
			_redis = client
		}
		e.Redis.Producer.RedisClient = client
		e.Redis.Consumer.RedisClient = client
		return queue.NewRedis(e.Redis.Producer, e.Redis.Consumer)
	}
	return queue.NewMemory(e.Memory.PoolSize), nil
}
