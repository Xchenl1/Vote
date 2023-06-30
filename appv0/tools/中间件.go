package tools

import (
	"book_manage_system/appv0/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 用户中间件
func UserAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		data, err := Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, HttpCode{
				Code:    NotLogin,
				Message: "验签失败！",
			})
		}
		if data.ID <= 0 {
			//如果用户没有登录，中间件直接返回，不再向后继续
			c.AbortWithStatusJSON(http.StatusUnauthorized, HttpCode{
				Code:    NotLogin,
				Message: "用户信息获取错误",
			})
			return
		}
		c.Set("userId", data.ID)
		c.Next()
		return
	}
}

// 管理员中间件
func AdminAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := GetSession(c)
		id, ok1 := data["id"]
		name, ok2 := data["name"]
		idInt64, _ := id.(int64)
		if !ok1 || !ok2 || idInt64 <= 0 || name == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() //如果用户没有登录，中间件直接返回，不再向后继续
		}
		c.Set("name", name)
		c.Set("userId", idInt64)
		c.Next()
	}
}

// 限制用户在一个端对单独的一个接口的访问
func LimitedFlow(maxCount int, t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取用户ip
		ip := c.ClientIP()
		//获取用户UA
		ua := c.GetHeader("User-Agent")
		//获取用户访问资源路径和请求参数，
		url := strings.Split(c.Request.URL.Path+"?"+c.Request.URL.RawQuery, "/")
		lastUrl := url[len(url)-1:]
		pathStr := fmt.Sprintf("%v_%v_%v", ip, ua, lastUrl)
		//fmt.Println("用户访问:", pathStr)
		requestQuery := model.RedisConn
		//访问的路径次数加1
		requestQuery.Incr(c, pathStr)
		//获取访问次数
		reqCountString, _ := requestQuery.Get(c, pathStr).Result()
		reqCount, _ := strconv.Atoi(reqCountString)
		//超过最大次数，限制访问
		if reqCount > maxCount {
			c.JSON(http.StatusOK, HttpCode{
				Code:    OK,
				Message: "请求太快，请休息一下重试~",
				Data:    nil,
			})
			c.Abort()
			return
		} else if reqCount == 1 { //第一次访问，设置过期时间为t
			// 设置键 "counter" 的过期时间为 120 秒
			if _, err := requestQuery.Expire(c, pathStr, t).Result(); err != nil {
				c.JSON(http.StatusOK, HttpCode{
					Code:    OK,
					Message: "未知错误",
				})
				c.Abort()
				return
			}
		}
		//每次请求时，将请求url存到redis
		c.Next()
	}
}
