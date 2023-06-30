package tools

import (
	"book_manage_system/appv0/model"
	"time"
)

func TimerMiddleware() {
	// 定义一个定时器，每10秒钟触发一次
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	//var mutex sync.Mutex
	// 在goroutine中运行定时器
	go func() {
		for {
			select {
			case <-ticker.C:
				//mutex.Lock()
				mesg := model.Sendmesg()
				model.Cr(mesg)
				//mutex.Unlock()
			}
		}
	}()
	//return func(c *gin.Context) {
	//	// 中间件处理逻辑
	//	c.Next()
	//}
}
