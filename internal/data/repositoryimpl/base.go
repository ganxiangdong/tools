package repositoryimpl

import (
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/internal/biz/bo"
	"github.com/mitchellh/mapstructure"
)

type Query[Q any] interface {
	Limit(limit int) *Q
	Offset(offset int) *Q
}

// Base 基础类型，提供了 po与do之间的转换等共用方法，在do中嵌入此类型创建，大部分情况下可不用再重复写 ToEntity 和 ToEntities 方法
type Base[P, D, Q any] struct {
}

// ToEntity 转换成实体
// 支持基本类型的值对象
// 支持 edge 的实体转换，注意名称要保持一致，否则只能在组合类覆盖此此方法
func (Base[P, D, Q]) ToEntity(p *P) *D {
	if p == nil {
		return nil
	}
	var d *D
	_ = mapstructure.Decode(p, &d)
	return d
}

// ToEntities 转换成实体
// 支持基本类型的值对象
// 支持 edge 的实体转换，注意名称要保持一致，否则只能在组合类覆盖此此方法
func (b Base[P, D, Q]) ToEntities(ps []*P) []*D {
	if ps == nil {
		return nil
	}
	entities := make([]*D, len(ps))
	for k, p := range ps {
		entities[k] = b.ToEntity(p)
	}
	return entities
}

// SetPageByBo 设置分页
func (b Base[P, D, Q]) SetPageByBo(query Query[Q], pageBo *bo.ReqPageBo) {
	if query == nil || pageBo == nil {
		return
	}
	if pageBo.Size > 0 {
		query.Limit(pageBo.Size)
	}
	if pageBo.Num > 0 {
		query.Offset(pageBo.GetOffset())
	}
}
