package dto

// PaginationRequest 通用分页请求参数
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`         // 页码，从1开始
	PageSize int `form:"pageSize" json:"pageSize" binding:"omitempty,min=1"` // 每页数量
}

// GetPage 获取页码（提供默认值）
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量（提供默认值）
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	// 限制最大值
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// GetLimit 获取限制数量
func (p *PaginationRequest) GetLimit() int {
	return p.GetPageSize()
}

// SortOrder 排序方向
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SortRequest 排序请求参数
type SortRequest struct {
	SortBy    string    `form:"sortBy" json:"sortBy"`       // 排序字段
	SortOrder SortOrder `form:"sortOrder" json:"sortOrder"` // 排序方向：asc, desc
}

// GetSortBy 获取排序字段（提供默认值）
func (s *SortRequest) GetSortBy(defaultField string) string {
	if s.SortBy == "" {
		return defaultField
	}
	return s.SortBy
}

// GetSortOrder 获取排序方向（提供默认值）
func (s *SortRequest) GetSortOrder() SortOrder {
	if s.SortOrder == "" {
		return SortOrderDesc
	}
	return s.SortOrder
}

// PageSortRequest 分页+排序请求参数
type PageSortRequest struct {
	PaginationRequest
	SortRequest
}
