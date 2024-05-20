package repository

import (
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/internal/biz/bo"
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/internal/biz/do"
	"context"
)

type MemberRepo interface {
	// Get 通过 id 获取一条数据，出错则 panic
	Get(ctx context.Context, id int) *do.MemberDo

	// Find 通过多个 id 获取多条数据，出错则 panic
	Find(ctx context.Context, ids ...int) []*do.MemberDo

	// Create 创建数据，出错则 panic
	Create(ctx context.Context, d *do.MemberDo) *do.MemberDo

	// CreateBulk 批量创建数据，出错则 panic
	CreateBulk(ctx context.Context, dos []*do.MemberDo) []*do.MemberDo

	// Update 更新数据，出错则 panic
	Update(ctx context.Context, d *do.MemberDo) int

	// Delete 删除数据，出错则 panic
	Delete(ctx context.Context, ids ...int) int

	// SearchList 搜索列表，出错则 panic
	SearchList(ctx context.Context, reqBo *bo.ReqPageBo) []*do.MemberDo
}
