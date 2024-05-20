//go:build testdata

package application

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

const OperationUserGetOptions = "/apiv1.User/GetOptions"

type UserHTTPServer interface {
	// GetOptions 获取用户下拉选项
	GetOptions(context.Context, *ReqUserGetOption) (*RespIntOptions, error)
}

func RegisterUserHTTPServer(s *http.Server, srv UserHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/user/options", _User_GetOptions4_HTTP_Handler(srv))
}

func _User_GetOptions4_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ReqUserGetOption
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserGetOptions)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetOptions(ctx, req.(*ReqUserGetOption))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RespIntOptions)
		return ctx.Result(200, reply)
	}
}

type UserHTTPClient interface {
	GetOptions(ctx context.Context, req *ReqUserGetOption, opts ...http.CallOption) (rsp *RespIntOptions, err error)
}

type UserHTTPClientImpl struct {
	cc *http.Client
}

func NewUserHTTPClient(client *http.Client) UserHTTPClient {
	return &UserHTTPClientImpl{client}
}

func (c *UserHTTPClientImpl) GetOptions(ctx context.Context, in *ReqUserGetOption, opts ...http.CallOption) (*RespIntOptions, error) {
	var out RespIntOptions
	pattern := "/v1/user/options"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserGetOptions))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
