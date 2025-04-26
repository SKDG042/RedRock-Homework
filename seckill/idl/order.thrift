namespace go order

// 基础的response
struct BaseResponse{
    1: i32      code     // 返回响应的状态码，0表示成功
    2: string   msg      // 返回响应信息
}

enum OrderStatus{
    PENDING   = 0,        // 订单创建中/订单未支付
    CREATED   = 1,        // 订单创建成功
    PAID      = 2,        // 订单已支付
    FAILED    = 3,        // 订单创建失败
    CANCELLED = 4,        // 订单已取消
}

// 订单信息
struct OrderInfo{
    1: i64          id          // 订单ID
    2: string       orderSn     // 订单号
    3: i64          userID      // 用户ID
    4: i64          activityID  // 活动ID
    5: i64          productID   // 商品ID
    6: string       ProductName // 商品名称
    // 作业文档中提到默认只购买一个商品，所以这里没有设置数量字段
    7: double       amount      // 订单总金额
    8: OrderStatus  status      // 订单状态
    9: i64          createTime  // 订单创建时间戳
    10:i64          payTime     // 订单支付时间戳
    11:i64          expireTime  // 订单过期时间戳
}

// 创建订单
struct CreateOrderRequest{
    1: i64          userID      // 用户ID
    2: i64          activityID  // 活动ID
}

struct CreateOrderResponse{
    1: BaseResponse baseResponse
    2: OrderInfo    orderInfo   // 订单信息
}

struct GetOrderRequest{
    1: i64          userID      // 用户ID
    2: string       orderSn     // 订单号
}

struct GetOrderResponse{
    1: BaseResponse baseResponse
    2: OrderInfo    orderInfo   // 订单信息
}

// 获取用户的订单列表
struct ListOrdersRequest{
    1: i64          userID      // 用户ID
    2: OrderStatus  status = -1 // 订单状态，-1表示所有订单
}

struct ListOrdersResponse{
    1: BaseResponse baseResponse
    2: list<OrderInfo> orders   // 订单列表
    3: i64              total   // 订单总数
}

service OrderService{
    // 创建订单
    CreateOrderResponse CreateOrder(1:CreateOrderRequest req)

    // 获取订单
    GetOrderResponse GetOrder(1:GetOrderRequest req)

    // 获取用户所有订单
    ListOrdersResponse ListOrders(1:ListOrdersRequest req)
}
