package main

// TestOrderStatus 订单状态
//
//go:generate stringer -type=TestOrderStatus -linecomment
type TestOrderStatus int

const (
	TestOrderStatusPaid   TestOrderStatus = iota + 1 // 已支付
	TestOrderStatusRefund                            // 已退款
	TestOrderStatusCancel                            // 已取消
)

// idea中调试，在program params中添加：-type=TestOrderStatus -linecomment 测试文件的绝对目录[如/Users/gangan/www/lsxd/tools/cmd/stringer]
