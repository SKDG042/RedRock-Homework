package data

import(
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	redisClient "Redrock/seckill/internal/pkg/redis"
	"Redrock/seckill/internal/pkg/models"
)

const(
	// 活动信息缓存键名前缀
	activityCacheKeyPrefix = "activity:info:"

	// 活动库存缓存键名前缀
	StockCacheKeyPrefix = "activity:stock:"

	// 用户参与活动记录简明前缀
	userJoinKeyPrefix = "activity:join:user:"

	// 缓存数据过期时间
	cacheExpireTime = 24 * time.Hour // 默认过期时间为24小时
)

type ActivityRedis struct{
	client *redis.Client
}

// NewActtivityRedis 创建Redis操作对象
func NewActivityRedis() *ActivityRedis{
	return &ActivityRedis{
		client : redisClient.GetRedis(),
	}
}

// SaveActivity 将活动信息保存到Redis
func (r *ActivityRedis) SaveActivity(ctx context.Context, activity *models.Activity) error{
	key :=  fmt.Sprintf("%s%d", activityCacheKeyPrefix, activity.ID)

	// 将数据序列化为json后保存在redis
	// func json.Marshal(v any) ([]byte, error)
	data, err := json.Marshal(activity)
	if err != nil{
		return err
	}
	// client.set()需要将字节切片转化为字符串
	err = r.client.Set(ctx, key, string(data), cacheExpireTime).Err()
	return err
}

// GetActivity 从Redis中获取活动信息
func (r *ActivityRedis) GetActivity(ctx context.Context, id uint) (*models.Activity, error){
	key := fmt.Sprintf("%s%d", activityCacheKeyPrefix, id)

	// 从redis中获取活动信息
	data, err := r.client.Get(ctx, key).Result()
	if err != nil{
		if err == redis.Nil{
			// 缓存不存在
			return nil, nil
		}
		return nil, err
	}

	// 因为从redis获取的是json格式字符串，需要反序列化
	var activity models.Activity
	err = json.Unmarshal([]byte(data), &activity)
	if err != nil{
		return nil,err
	}
	
	return &activity, nil
}

// InitStock 初始化库存信息到Redis
func (r *ActivityRedis) InitStock(ctx context.Context, activityID uint, stock int64) (int64, error){
	key := fmt.Sprintf("%s%d", StockCacheKeyPrefix, activityID)

	// 使用SetNX而不是Set的原因：
	// 仅在键不存在时设置key的值，即初始化
	err := r.client.SetNX(ctx, key, stock, cacheExpireTime).Err()

	return stock,err
}

// GetStock 获取当前库存
func (r *ActivityRedis) GetStock(ctx context.Context, activityID uint) (int64, error){
	key := fmt.Sprintf("%s%d", StockCacheKeyPrefix, activityID)

	// 获取当前库存 （返回string）
	stockStr, err := r.client.Get(ctx, key).Result()
	if err != nil{
		if err == redis.Nil{
			// 缓存不存在
			return 0, nil
		}
		return 0, err
	}

	stock, err := strconv.ParseInt(stockStr, 10, 64)
	if err != nil{
		return 0, err
	}

	return stock, nil
}

// DeductStock 扣除库存 (Redis的事务无法回退，因此使用lua脚本保证原子性)
func(r *ActivityRedis) DeductStock(ctx context.Context, activityID uint, count int64) (bool, error){
	key := fmt.Sprintf("%s%d", StockCacheKeyPrefix, activityID)

	// 使用lua脚本扣除库存以保证原子性
	script := `
	local stock = tonumber(redis.call("GET", KEYS[1]))
	if stock == nil then
		return -1 -- 库存不存在
	end

	if stock < tonumber(ARGV[1]) then
		return 0 -- 库存不足
	end

	redis.call("decrby", KEYS[1], ARGV[1])
		return 1 -- 扣除成功
	`
	// 执行lua脚本
	result, err := r.client.Eval(ctx, script, []string{key}, count).Int64()
	if err != nil{
		return false, err
	}

	if result == -1{
		return false, fmt.Errorf("库存信息不存在")
	}

	return result == 1, nil
}

/// RecordUserJoin 记录用户参与记录
func (r *ActivityRedis) RecordUserJoin(ctx context.Context, userID uint, activityID uint) error{
	key := fmt.Sprintf("%s%d:%d",userJoinKeyPrefix,userID,activityID)

	// 设置用户参与记录以及有效期
	err := r.client.SetNX(ctx, key, 1, cacheExpireTime).Err()
	return err
}

// IsUserJoined 检查用户是否已参与活动
func(r *ActivityRedis) IsUserJoined(ctx context.Context, userID uint, activityID uint) (bool, error){
	key := fmt.Sprintf("%s%d:%d", userJoinKeyPrefix,userID,activityID)

	// 检查记录是否存在
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil{
		return false, err
	}

	return exists == 1, nil
}
