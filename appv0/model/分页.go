package model

import (
	"strconv"
)

type Page[T any] struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	Total       int `json:"total"` //总数
	Pages       int `json:"pages"` //总页数
	Result      []T `json:"result"`
}

// Pages 这部分代码是需要输入当前页码 以及每一页的大小
func Pages[T any](res []T, currentPageString, pageSizeString string) Page[T] {
	currentPage, _ := strconv.Atoi(currentPageString)
	pageSize, _ := strconv.Atoi(pageSizeString)
	offset := (currentPage - 1) * pageSize
	limit := pageSize
	result := res[offset : offset+limit]
	if len(result) == 0 {
		return Page[T]{}
	}
	page := Page[T]{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		Total:       len(res),
		Pages:       len(res)/pageSize + 1,
		Result:      result,
	}
	return page
}

// BatchGet 分批获取书籍信息可以将较大的数据集分成若干批次进行处理，
// 从而避免一次性获取过多数据导致内存不足或响应时间过长的问题
func BatchGet[T any](res []T, batchSize int) []Page[T] {
	var pages []Page[T]
	currentPage := 1
	pageSize := batchSize
	total := len(res)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		result := res[i:end]
		page := Page[T]{
			CurrentPage: currentPage,
			PageSize:    pageSize,
			Total:       total,
			Pages:       (total + batchSize - 1) / batchSize,
			Result:      result,
		}
		pages = append(pages, page)
		currentPage++
	}
	return pages
}
