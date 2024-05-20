package biz

// TestABiz 用于测试 wire_cmd 命令
type TestABiz struct {
}

func NewTestABiz() *TestABiz {
	return &TestABiz{}
}

func (t *TestABiz) Test() *TestABiz {
	return &TestABiz{}
}
