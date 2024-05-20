package biz

import "github.com/google/wire"

var ProviderSetService = wire.NewSet(
	NewTestBBiz,
	NewTestABiz,
)
