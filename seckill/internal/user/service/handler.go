package service

import (
	"context"

	user "Redrock/seckill/kitex_gen/user"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/user/data"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{
	userData *data.UserData
}

func NewUserServiceImpl() *UserServiceImpl{
	return &UserServiceImpl{
		userData: data.NewUserData(),
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	response := &user.RegisterResponse{
		BaseResp: &user.BaseResp{},
	}

	// 参数验证
	if req.Username == "" || req.Password == "" {
		response.BaseResp.Code = 400
		response.BaseResp.Message = "用户名和密码不能为空"
		return response, nil
	}

	// 创建用户
	newUser := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	err = s.userData.Create(ctx, newUser)
	if err != nil {
		response.BaseResp.Code = 500
		response.BaseResp.Message = "创建用户失败: " + err.Error()
		return response, nil
	}

	// 返回成功
	response.BaseResp.Code = 0
	response.BaseResp.Message = "注册成功"
	response.UserId = int64(newUser.ID)

	return response, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	response := &user.LoginResponse{
		BaseResp: &user.BaseResp{},
	}

	// 参数验证
	if req.Username == "" || req.Password == "" {
		response.BaseResp.Code = 400
		response.BaseResp.Message = "用户名和密码不能为空"
		return response, nil
	}

	// 检查用户名密码
	loginUser, err := s.userData.CheckPassword(ctx, req.Username, req.Password)
	if err != nil {
		response.BaseResp.Code = 401
		response.BaseResp.Message = "登录失败: " + err.Error()
		return response, nil
	}

	response.BaseResp.Code = 0
	response.BaseResp.Message = "登录成功"
	response.UserId = int64(loginUser.ID)

	return response, nil
}
