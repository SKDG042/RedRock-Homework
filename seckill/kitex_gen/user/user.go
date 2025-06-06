// Code generated by thriftgo (0.4.1). DO NOT EDIT.

package user

import (
	"context"
	"fmt"
)

type BaseResp struct {
	Code    int32  `thrift:"code,1" frugal:"1,default,i32" json:"code"`
	Message string `thrift:"message,2" frugal:"2,default,string" json:"message"`
}

func NewBaseResp() *BaseResp {
	return &BaseResp{}
}

func (p *BaseResp) InitDefault() {
}

func (p *BaseResp) GetCode() (v int32) {
	return p.Code
}

func (p *BaseResp) GetMessage() (v string) {
	return p.Message
}
func (p *BaseResp) SetCode(val int32) {
	p.Code = val
}
func (p *BaseResp) SetMessage(val string) {
	p.Message = val
}

func (p *BaseResp) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseResp(%+v)", *p)
}

var fieldIDToName_BaseResp = map[int16]string{
	1: "code",
	2: "message",
}

type UserInfo struct {
	Id       int64  `thrift:"id,1" frugal:"1,default,i64" json:"id"`
	Username string `thrift:"username,2" frugal:"2,default,string" json:"username"`
}

func NewUserInfo() *UserInfo {
	return &UserInfo{}
}

func (p *UserInfo) InitDefault() {
}

func (p *UserInfo) GetId() (v int64) {
	return p.Id
}

func (p *UserInfo) GetUsername() (v string) {
	return p.Username
}
func (p *UserInfo) SetId(val int64) {
	p.Id = val
}
func (p *UserInfo) SetUsername(val string) {
	p.Username = val
}

func (p *UserInfo) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserInfo(%+v)", *p)
}

var fieldIDToName_UserInfo = map[int16]string{
	1: "id",
	2: "username",
}

type RegisterRequest struct {
	Username string `thrift:"username,1,required" frugal:"1,required,string" json:"username"`
	Password string `thrift:"password,2,required" frugal:"2,required,string" json:"password"`
}

func NewRegisterRequest() *RegisterRequest {
	return &RegisterRequest{}
}

func (p *RegisterRequest) InitDefault() {
}

func (p *RegisterRequest) GetUsername() (v string) {
	return p.Username
}

func (p *RegisterRequest) GetPassword() (v string) {
	return p.Password
}
func (p *RegisterRequest) SetUsername(val string) {
	p.Username = val
}
func (p *RegisterRequest) SetPassword(val string) {
	p.Password = val
}

func (p *RegisterRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RegisterRequest(%+v)", *p)
}

var fieldIDToName_RegisterRequest = map[int16]string{
	1: "username",
	2: "password",
}

type RegisterResponse struct {
	BaseResp *BaseResp `thrift:"baseResp,1" frugal:"1,default,BaseResp" json:"baseResp"`
	UserId   int64     `thrift:"userId,2" frugal:"2,default,i64" json:"userId"`
}

func NewRegisterResponse() *RegisterResponse {
	return &RegisterResponse{}
}

func (p *RegisterResponse) InitDefault() {
}

var RegisterResponse_BaseResp_DEFAULT *BaseResp

func (p *RegisterResponse) GetBaseResp() (v *BaseResp) {
	if !p.IsSetBaseResp() {
		return RegisterResponse_BaseResp_DEFAULT
	}
	return p.BaseResp
}

func (p *RegisterResponse) GetUserId() (v int64) {
	return p.UserId
}
func (p *RegisterResponse) SetBaseResp(val *BaseResp) {
	p.BaseResp = val
}
func (p *RegisterResponse) SetUserId(val int64) {
	p.UserId = val
}

func (p *RegisterResponse) IsSetBaseResp() bool {
	return p.BaseResp != nil
}

func (p *RegisterResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RegisterResponse(%+v)", *p)
}

var fieldIDToName_RegisterResponse = map[int16]string{
	1: "baseResp",
	2: "userId",
}

type LoginRequest struct {
	Username string `thrift:"username,1,required" frugal:"1,required,string" json:"username"`
	Password string `thrift:"password,2,required" frugal:"2,required,string" json:"password"`
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

func (p *LoginRequest) InitDefault() {
}

func (p *LoginRequest) GetUsername() (v string) {
	return p.Username
}

func (p *LoginRequest) GetPassword() (v string) {
	return p.Password
}
func (p *LoginRequest) SetUsername(val string) {
	p.Username = val
}
func (p *LoginRequest) SetPassword(val string) {
	p.Password = val
}

func (p *LoginRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LoginRequest(%+v)", *p)
}

var fieldIDToName_LoginRequest = map[int16]string{
	1: "username",
	2: "password",
}

type LoginResponse struct {
	BaseResp *BaseResp `thrift:"baseResp,1" frugal:"1,default,BaseResp" json:"baseResp"`
	UserId   int64     `thrift:"userId,2" frugal:"2,default,i64" json:"userId"`
}

func NewLoginResponse() *LoginResponse {
	return &LoginResponse{}
}

func (p *LoginResponse) InitDefault() {
}

var LoginResponse_BaseResp_DEFAULT *BaseResp

func (p *LoginResponse) GetBaseResp() (v *BaseResp) {
	if !p.IsSetBaseResp() {
		return LoginResponse_BaseResp_DEFAULT
	}
	return p.BaseResp
}

func (p *LoginResponse) GetUserId() (v int64) {
	return p.UserId
}
func (p *LoginResponse) SetBaseResp(val *BaseResp) {
	p.BaseResp = val
}
func (p *LoginResponse) SetUserId(val int64) {
	p.UserId = val
}

func (p *LoginResponse) IsSetBaseResp() bool {
	return p.BaseResp != nil
}

func (p *LoginResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LoginResponse(%+v)", *p)
}

var fieldIDToName_LoginResponse = map[int16]string{
	1: "baseResp",
	2: "userId",
}

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (r *RegisterResponse, err error)

	Login(ctx context.Context, req *LoginRequest) (r *LoginResponse, err error)
}

type UserServiceRegisterArgs struct {
	Req *RegisterRequest `thrift:"req,1" frugal:"1,default,RegisterRequest" json:"req"`
}

func NewUserServiceRegisterArgs() *UserServiceRegisterArgs {
	return &UserServiceRegisterArgs{}
}

func (p *UserServiceRegisterArgs) InitDefault() {
}

var UserServiceRegisterArgs_Req_DEFAULT *RegisterRequest

func (p *UserServiceRegisterArgs) GetReq() (v *RegisterRequest) {
	if !p.IsSetReq() {
		return UserServiceRegisterArgs_Req_DEFAULT
	}
	return p.Req
}
func (p *UserServiceRegisterArgs) SetReq(val *RegisterRequest) {
	p.Req = val
}

func (p *UserServiceRegisterArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *UserServiceRegisterArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserServiceRegisterArgs(%+v)", *p)
}

var fieldIDToName_UserServiceRegisterArgs = map[int16]string{
	1: "req",
}

type UserServiceRegisterResult struct {
	Success *RegisterResponse `thrift:"success,0,optional" frugal:"0,optional,RegisterResponse" json:"success,omitempty"`
}

func NewUserServiceRegisterResult() *UserServiceRegisterResult {
	return &UserServiceRegisterResult{}
}

func (p *UserServiceRegisterResult) InitDefault() {
}

var UserServiceRegisterResult_Success_DEFAULT *RegisterResponse

func (p *UserServiceRegisterResult) GetSuccess() (v *RegisterResponse) {
	if !p.IsSetSuccess() {
		return UserServiceRegisterResult_Success_DEFAULT
	}
	return p.Success
}
func (p *UserServiceRegisterResult) SetSuccess(x interface{}) {
	p.Success = x.(*RegisterResponse)
}

func (p *UserServiceRegisterResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *UserServiceRegisterResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserServiceRegisterResult(%+v)", *p)
}

var fieldIDToName_UserServiceRegisterResult = map[int16]string{
	0: "success",
}

type UserServiceLoginArgs struct {
	Req *LoginRequest `thrift:"req,1" frugal:"1,default,LoginRequest" json:"req"`
}

func NewUserServiceLoginArgs() *UserServiceLoginArgs {
	return &UserServiceLoginArgs{}
}

func (p *UserServiceLoginArgs) InitDefault() {
}

var UserServiceLoginArgs_Req_DEFAULT *LoginRequest

func (p *UserServiceLoginArgs) GetReq() (v *LoginRequest) {
	if !p.IsSetReq() {
		return UserServiceLoginArgs_Req_DEFAULT
	}
	return p.Req
}
func (p *UserServiceLoginArgs) SetReq(val *LoginRequest) {
	p.Req = val
}

func (p *UserServiceLoginArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *UserServiceLoginArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserServiceLoginArgs(%+v)", *p)
}

var fieldIDToName_UserServiceLoginArgs = map[int16]string{
	1: "req",
}

type UserServiceLoginResult struct {
	Success *LoginResponse `thrift:"success,0,optional" frugal:"0,optional,LoginResponse" json:"success,omitempty"`
}

func NewUserServiceLoginResult() *UserServiceLoginResult {
	return &UserServiceLoginResult{}
}

func (p *UserServiceLoginResult) InitDefault() {
}

var UserServiceLoginResult_Success_DEFAULT *LoginResponse

func (p *UserServiceLoginResult) GetSuccess() (v *LoginResponse) {
	if !p.IsSetSuccess() {
		return UserServiceLoginResult_Success_DEFAULT
	}
	return p.Success
}
func (p *UserServiceLoginResult) SetSuccess(x interface{}) {
	p.Success = x.(*LoginResponse)
}

func (p *UserServiceLoginResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *UserServiceLoginResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserServiceLoginResult(%+v)", *p)
}

var fieldIDToName_UserServiceLoginResult = map[int16]string{
	0: "success",
}
