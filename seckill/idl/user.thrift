namespace go user

// 基础响应结构
struct BaseResp {
    1: i32 code                // 状态码，0表示成功，非0表示各种错误
    2: string message          // 错误信息
}

// 用户信息
struct UserInfo {
    1: i64 id                  // 用户ID
    2: string username         // 用户名
}

// 用户注册请求
struct RegisterRequest {
    1: required string username    // 用户名
    2: required string password    // 密码
}

// 用户注册响应
struct RegisterResponse {
    1: BaseResp baseResp       // 基础响应
    2: i64 userId              // 用户ID
}

// 用户登录请求
struct LoginRequest {
    1: required string username    // 用户名
    2: required string password    // 密码
}

// 用户登录响应
struct LoginResponse {
    1: BaseResp baseResp       // 基础响应
    2: i64 userId              // 用户ID
}

// 用户服务
service UserService {
    // 用户注册
    RegisterResponse Register(1: RegisterRequest req)
    
    // 用户登录
    LoginResponse Login(1: LoginRequest req)
}
