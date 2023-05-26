// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             v3.21.12
// source: admin/v1/api.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAppCategory = "/api.admin.v1.App/Category"
const OperationAppGetHost = "/api.admin.v1.App/GetHost"
const OperationAppLogin = "/api.admin.v1.App/Login"

type AppHTTPServer interface {
	Category(context.Context, *CategoryReq) (*CategoryResp, error)
	GetHost(context.Context, *GetHostReq) (*GetHostResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
}

func RegisterAppHTTPServer(s *http.Server, srv AppHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/login", _App_Login0_HTTP_Handler(srv))
	r.GET("/v1/category", _App_Category0_HTTP_Handler(srv))
	r.GET("/v1/get_host", _App_GetHost0_HTTP_Handler(srv))
}

func _App_Login0_HTTP_Handler(srv AppHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAppLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResp)
		return ctx.Result(200, reply)
	}
}

func _App_Category0_HTTP_Handler(srv AppHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CategoryReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAppCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Category(ctx, req.(*CategoryReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CategoryResp)
		return ctx.Result(200, reply)
	}
}

func _App_GetHost0_HTTP_Handler(srv AppHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetHostReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAppGetHost)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetHost(ctx, req.(*GetHostReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetHostResp)
		return ctx.Result(200, reply)
	}
}

type AppHTTPClient interface {
	Category(ctx context.Context, req *CategoryReq, opts ...http.CallOption) (rsp *CategoryResp, err error)
	GetHost(ctx context.Context, req *GetHostReq, opts ...http.CallOption) (rsp *GetHostResp, err error)
	Login(ctx context.Context, req *LoginReq, opts ...http.CallOption) (rsp *LoginResp, err error)
}

type AppHTTPClientImpl struct {
	cc *http.Client
}

func NewAppHTTPClient(client *http.Client) AppHTTPClient {
	return &AppHTTPClientImpl{client}
}

func (c *AppHTTPClientImpl) Category(ctx context.Context, in *CategoryReq, opts ...http.CallOption) (*CategoryResp, error) {
	var out CategoryResp
	pattern := "/v1/category"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAppCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AppHTTPClientImpl) GetHost(ctx context.Context, in *GetHostReq, opts ...http.CallOption) (*GetHostResp, error) {
	var out GetHostResp
	pattern := "/v1/get_host"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAppGetHost))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AppHTTPClientImpl) Login(ctx context.Context, in *LoginReq, opts ...http.CallOption) (*LoginResp, error) {
	var out LoginResp
	pattern := "/v1/login"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAppLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
