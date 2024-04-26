// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user.proto

package user

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for User service

func NewUserEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for User service

type UserService interface {
	SendSms(ctx context.Context, in *SmsRequest, opts ...client.CallOption) (*SmsResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) SendSms(ctx context.Context, in *SmsRequest, opts ...client.CallOption) (*SmsResponse, error) {
	req := c.c.NewRequest(c.name, "User.SendSms", in)
	out := new(SmsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for User service

type UserHandler interface {
	SendSms(context.Context, *SmsRequest, *SmsResponse) error
}

func RegisterUserHandler(s server.Server, hdlr UserHandler, opts ...server.HandlerOption) error {
	type user interface {
		SendSms(ctx context.Context, in *SmsRequest, out *SmsResponse) error
	}
	type User struct {
		user
	}
	h := &userHandler{hdlr}
	return s.Handle(s.NewHandler(&User{h}, opts...))
}

type userHandler struct {
	UserHandler
}

func (h *userHandler) SendSms(ctx context.Context, in *SmsRequest, out *SmsResponse) error {
	return h.UserHandler.SendSms(ctx, in, out)
}
