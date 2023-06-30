package tools

import (
	"fmt"
	sessions2 "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"time"
)

// 导入session包 命令：go get github.com/gorilla/sessions
var store, _ = redis.NewStore(10, "tcp", "localhost:6379", "", []byte("小陈图书管理系统"))

// var store = sessions.NewCookieStore([]byte("小陈图书管理系统"))
var sessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(c *gin.Context, name string, id int64) error {
	//配置session中redis中的生存周期24小时的生存周期
	store.Options(sessions2.Options(sessions.Options{
		MaxAge: int(24 * time.Hour / time.Second),
	}))
	session, _ := store.Get(c.Request, sessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = 0
	return session.Save(c.Request, c.Writer)
}
