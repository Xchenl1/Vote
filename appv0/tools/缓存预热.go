package tools

import (
	"book_manage_system/appv0/model"
	"time"
)

// Cacheheating  缓存预热
func Cacheheating() {
	// 创建一个定时器，每 3秒触发一次
	ticker1 := time.NewTicker(3 * time.Second)
	defer ticker1.Stop()
	// 设置初始页码为 1
	id := 1
	size := 100
	// 在goroutine中运行定时器
	for {
		select {
		case <-ticker1.C:
			// 调用 Loader 函数加载图书信息
			books := model.Loader(id, size)
			// 调用 Handler 函数处理图书信息
			model.Handler(id, books)
			// 修改 pageIndex 变量，准备加载下一页图书信息
			id += size
			if id == 301 {
				id = 1
			}
		}
	}
}
