package bo

// ReqPageBo 分页请求实体
type ReqPageBo struct {
	Size int //分页大小
	Num  int //页码，从第1页开始
}

// GetOffset 获取便宜量
func (r *ReqPageBo) GetOffset() int {
	if r == nil {
		return 0
	}
	offset := (r.Num - 1) * r.Size
	if offset < 0 {
		return 0
	}
	return offset
}
