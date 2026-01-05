package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

// Filter 查询过滤器接口
type Filter interface {
	Apply(query *bun.SelectQuery) *bun.SelectQuery
}

// WhereFilter 简单的WHERE条件过滤器
type WhereFilter struct {
	Column   string
	Operator string
	Value    interface{}
}

func (f *WhereFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	condition := fmt.Sprintf("%s %s ?", f.Column, f.Operator)
	return query.Where(condition, f.Value)
}

// NewWhereFilter 创建WHERE过滤器
func NewWhereFilter(column, operator string, value interface{}) *WhereFilter {
	return &WhereFilter{
		Column:   column,
		Operator: operator,
		Value:    value,
	}
}

// InFilter IN条件过滤器
type InFilter struct {
	Column string
	Values []interface{}
}

func (f *InFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	if len(f.Values) == 0 {
		return query
	}
	return query.Where(fmt.Sprintf("%s IN (?)", f.Column), bun.In(f.Values))
}

// NewInFilter 创建IN过滤器
func NewInFilter(column string, values []interface{}) *InFilter {
	return &InFilter{
		Column: column,
		Values: values,
	}
}

// LikeFilter LIKE条件过滤器
type LikeFilter struct {
	Column string
	Value  string
}

func (f *LikeFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	if f.Value == "" {
		return query
	}
	return query.Where(fmt.Sprintf("%s LIKE ?", f.Column), "%"+f.Value+"%")
}

// NewLikeFilter 创建LIKE过滤器
func NewLikeFilter(column, value string) *LikeFilter {
	return &LikeFilter{
		Column: column,
		Value:  value,
	}
}

// DateRangeFilter 日期范围过滤器
type DateRangeFilter struct {
	Column    string
	StartTime *time.Time
	EndTime   *time.Time
}

func (f *DateRangeFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	if f.StartTime != nil {
		query = query.Where(fmt.Sprintf("%s >= ?", f.Column), f.StartTime)
	}
	if f.EndTime != nil {
		query = query.Where(fmt.Sprintf("%s <= ?", f.Column), f.EndTime)
	}
	return query
}

// NewDateRangeFilter 创建日期范围过滤器
func NewDateRangeFilter(column string, startTime, endTime *time.Time) *DateRangeFilter {
	return &DateRangeFilter{
		Column:    column,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// RangeFilter 数值范围过滤器
type RangeFilter struct {
	Column string
	Min    *int64
	Max    *int64
}

func (f *RangeFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	if f.Min != nil {
		query = query.Where(fmt.Sprintf("%s >= ?", f.Column), *f.Min)
	}
	if f.Max != nil {
		query = query.Where(fmt.Sprintf("%s <= ?", f.Column), *f.Max)
	}
	return query
}

// NewRangeFilter 创建数值范围过滤器
func NewRangeFilter(column string, min, max *int64) *RangeFilter {
	return &RangeFilter{
		Column: column,
		Min:    min,
		Max:    max,
	}
}

// OrderFilter 排序过滤器
type OrderFilter struct {
	Column string
	Desc   bool
}

func (f *OrderFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	if f.Desc {
		return query.Order(f.Column + " DESC")
	}
	return query.Order(f.Column + " ASC")
}

// NewOrderFilter 创建排序过滤器
func NewOrderFilter(column string, desc bool) *OrderFilter {
	return &OrderFilter{
		Column: column,
		Desc:   desc,
	}
}

// PaginationFilter 分页过滤器
type PaginationFilter struct {
	Page     int
	PageSize int
}

func (f *PaginationFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	offset := (f.Page - 1) * f.PageSize
	return query.Offset(offset).Limit(f.PageSize)
}

// NewPaginationFilter 创建分页过滤器
func NewPaginationFilter(page, pageSize int) *PaginationFilter {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return &PaginationFilter{
		Page:     page,
		PageSize: pageSize,
	}
}

// GetOffset 获取偏移量
func (f *PaginationFilter) GetOffset() int {
	return (f.Page - 1) * f.PageSize
}

// GetLimit 获取限制数量
func (f *PaginationFilter) GetLimit() int {
	return f.PageSize
}

// CompositeFilter 组合过滤器
type CompositeFilter struct {
	Filters []Filter
}

func (f *CompositeFilter) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	for _, filter := range f.Filters {
		query = filter.Apply(query)
	}
	return query
}

// NewCompositeFilter 创建组合过滤器
func NewCompositeFilter(filters ...Filter) *CompositeFilter {
	return &CompositeFilter{
		Filters: filters,
	}
}

// Add 添加过滤器
func (f *CompositeFilter) Add(filter Filter) {
	f.Filters = append(f.Filters, filter)
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	filters []Filter
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		filters: make([]Filter, 0),
	}
}

// Where 添加WHERE条件
func (qb *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
	qb.filters = append(qb.filters, NewWhereFilter(column, operator, value))
	return qb
}

// In 添加IN条件
func (qb *QueryBuilder) In(column string, values []interface{}) *QueryBuilder {
	qb.filters = append(qb.filters, NewInFilter(column, values))
	return qb
}

// Like 添加LIKE条件
func (qb *QueryBuilder) Like(column, value string) *QueryBuilder {
	if value != "" {
		qb.filters = append(qb.filters, NewLikeFilter(column, value))
	}
	return qb
}

// DateRange 添加日期范围条件
func (qb *QueryBuilder) DateRange(column string, startTime, endTime *time.Time) *QueryBuilder {
	qb.filters = append(qb.filters, NewDateRangeFilter(column, startTime, endTime))
	return qb
}

// Range 添加数值范围条件
func (qb *QueryBuilder) Range(column string, min, max *int64) *QueryBuilder {
	qb.filters = append(qb.filters, NewRangeFilter(column, min, max))
	return qb
}

// Order 添加排序
func (qb *QueryBuilder) Order(column string, desc bool) *QueryBuilder {
	qb.filters = append(qb.filters, NewOrderFilter(column, desc))
	return qb
}

// Paginate 添加分页
func (qb *QueryBuilder) Paginate(page, pageSize int) *QueryBuilder {
	qb.filters = append(qb.filters, NewPaginationFilter(page, pageSize))
	return qb
}

// Apply 应用所有过滤器
func (qb *QueryBuilder) Apply(query *bun.SelectQuery) *bun.SelectQuery {
	for _, filter := range qb.filters {
		query = filter.Apply(query)
	}
	return query
}

// SafeColumn 安全的列名（防止SQL注入）
func SafeColumn(column string, allowedColumns []string) string {
	column = strings.TrimSpace(column)
	for _, allowed := range allowedColumns {
		if column == allowed {
			return column
		}
	}
	return ""
}
