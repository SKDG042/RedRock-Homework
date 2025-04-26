namespace go activity

// 基础的response
struct BaseResponse{
    1: i32      code     // 返回响应的状态码，0表示成功
    2: string   msg      // 返回响应信息
}

// 活动信息
struct ActivityInfo{
    1: i64      id              // 活动ID
    2: string   name            // 活动名称
    3: i64      productId       // 商品ID
    4: string   productName     // 商品名称
    5: double   originalPrice   // 商品原价
    6: double   seckillPrice    // 秒杀价格
    7: i64      startTime       // 活动开始时间戳(秒)
    8: i64      endTime         // 活动结束时间戳(喵)
    9: i64      totalStock      // 总库存
    10: i64     availableStock  // 可用库存
    11: bool    isAvailable     // 活动是否可用
}

// 创建活动请求
struct CreateActivityRequest{
    1: string   name            // 活动名称
    2: i64      productID       // 商品ID
    3: double   seckillPrice    // 秒杀价格
    4: i64      startTime       // 活动开始时间戳
    5: i64      endTime         // 活动结束时间戳
    6: i64      totalStock      // 总库存
}

// 常见活动响应
struct CreateActivityResponse{
    1: BaseResponse baseResponse
    2: i64          activityID      // 活动ID
}

// 获取活动列表
struct GetActivityListRequest{

}

struct GetActivityListResponse{
    1: BaseResponse         baseResponse
    2: list<ActivityInfo>   activities  //  活动列表
    3: i64                  total       // 总活动数量
}

// 获取活动
struct GetActivityRequest{
    1: i64                  activityID  // 活动ID
}

struct GetActivityResponse{
    1: BaseResponse         baseResponse
    2: ActivityInfo         activity    // 活动信息
}

// 扣除库存
struct DeductStockRequest{
    1: i64                  activityID  // 活动ID
    2: i64                  userID      // 用户ID
    3: i64                  count = 1   // 扣除数量，default 1
}

struct DeductStockResponse{
    1: BaseResponse         baseResponse
    2: bool                 success     // 是否成功扣除数量
}

service ActivityService{
    // 创建活动
    CreateActivityResponse      CreateActivity(1: CreateActivityRequest req)

    // 获取所有活动
    GetActivityListResponse     GetActivityList(1: GetActivityListRequest req)

    // 获取活动详情
    GetActivityResponse         GetActivity(1: GetActivityRequest req)
}

service InternalActivityService{
    // 扣除库存
    DeductStockResponse         DeductStock(1: DeductStockRequest req)
}
