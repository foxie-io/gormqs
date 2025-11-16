package dto

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/foxie-io/gormqs"
)

var (
	_ gormqs.ManyWithCountResulter = (*PageParam[any])(nil)
)

type PageResult[Slice any] struct {
	Data    *Slice `json:"users,omitempty"`
	Total   *int64 `json:"count,omitempty"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	HasNext *bool  `json:"hasNext,omitempty"`
}

type PageParam[List any] struct {
	result *PageResult[List]
	Select string `json:"select"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

func NewPageParam[List any](r *http.Request) *PageParam[List] {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 10 {
		size = 10
	}

	_select := r.URL.Query().Get("select")
	hasCount, hasList := strings.Contains(_select, "count"), strings.Contains(_select, "list")

	if !hasCount && !hasList {
		hasCount = true
		hasList = true
	}

	result := &PageResult[List]{
		Page: page,
		Size: size,
	}

	if hasList {
		result.Data = new(List)
	}

	if hasCount {
		result.Total = new(int64)
	}

	return &PageParam[List]{
		result: result,
		Select: _select,
		Page:   page,
		Size:   size,
	}
}

func (r *PageParam[List]) Result() *PageResult[List] {
	return r.result
}

func (r *PageParam[List]) HasNext(rows int) bool {
	if r.result.Total != nil {
		fullSize := int64(r.Size) * int64(r.Page)
		return fullSize < *r.result.Total
	}

	return rows <= r.Size
}

func (r *PageParam[List]) DBOption() gormqs.Option {
	return gormqs.LimitAndOffset(r.Size, (r.Page-1)*r.Size)
}

func (r *PageParam[List]) QsList() any {
	if r.result.Data != nil {
		return r.result.Data
	}
	return nil
}

func (r *PageParam[List]) QsCount() *int64 {
	return r.result.Total
}
