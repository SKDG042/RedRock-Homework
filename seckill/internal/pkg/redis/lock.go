package redis

import(
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

// DistributedLock Redis分布式锁
type DistributedLock struct{
	redisClient *redis.Client
	key			string
	value		string
	expiration	time.Duration
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(redisClient *redis.Client, key string, expiration time.Duration) *DistributedLock {
    return &DistributedLock{
        redisClient: redisClient,
        key:         "lock:" + key,
        value:       uuid.New().String(), // 使用唯一uuid作为value，保证只能由持有者释放
        expiration:  expiration,
    }
}

// TryLock 尝试获取锁
func (l *DistributedLock) TryLock(ctx context.Context) (bool, error){
	// 如果锁不存在则设置
	return l.redisClient.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

// Unlock
func (l *DistributedLock) Unlock(ctx context.Context) error{
	const script =
	`
	if redis.call("GET", KEY[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end	
	`
	// 当key和value匹配时才释放，防止错释放他人的锁
	_, err := l.redisClient.Eval(ctx, script, []string{l.key}, l.value).Result()
	return err
}
