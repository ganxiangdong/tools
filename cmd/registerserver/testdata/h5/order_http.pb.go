//go:build testdata

package h5

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

const OperationOrderCreate = "/api.alipayoil.openapi.Order/Create"
const OperationOrderQuery = "/api.alipayoil.openapi.Order/Query"

type OrderHTTPServer interface {
	// Create 创建订单
	Create(context.Context, *OrderCreateReq) (*OrderCreateResp, error)
	// Query 查询订单
	Query(context.Context, *OrderQueryReq) (*OrderQueryResp, error)
}

func RegisterOrderHTTPServer(s *http.Server, srv OrderHTTPServer) {
	r := s.Route("/")
	r.POST("/openapi/order/create", _Order_Create0_HTTP_Handler(srv))
	r.POST("/openapi/order/query", _Order_Query0_HTTP_Handler(srv))
}

func _Order_Create0_HTTP_Handler(srv OrderHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OrderCreateReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOrderCreate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Create(ctx, req.(*OrderCreateReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OrderCreateResp)
		return ctx.Result(200, reply)
	}
}

func _Order_Query0_HTTP_Handler(srv OrderHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OrderQueryReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOrderQuery)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Query(ctx, req.(*OrderQueryReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OrderQueryResp)
		return ctx.Result(200, reply)
	}
}

type OrderHTTPClient interface {
	Create(ctx context.Context, req *OrderCreateReq, opts ...http.CallOption) (rsp *OrderCreateResp, err error)
	Query(ctx context.Context, req *OrderQueryReq, opts ...http.CallOption) (rsp *OrderQueryResp, err error)
}

type OrderHTTPClientImpl struct {
	cc *http.Client
}

func NewOrderHTTPClient(client *http.Client) OrderHTTPClient {
	return &OrderHTTPClientImpl{client}
}

func (c *OrderHTTPClientImpl) Create(ctx context.Context, in *OrderCreateReq, opts ...http.CallOption) (*OrderCreateResp, error) {
	var out OrderCreateResp
	pattern := "/openapi/order/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOrderCreate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *OrderHTTPClientImpl) Query(ctx context.Context, in *OrderQueryReq, opts ...http.CallOption) (*OrderQueryResp, error) {
	var out OrderQueryResp
	pattern := "/openapi/order/query"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOrderQuery))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
