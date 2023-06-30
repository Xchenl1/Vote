package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

// SelectBook 查询全部书籍
func SelectBook(c *gin.Context, idstr string, id int, size int) []Book {
	// 从Redis中获取分页数据
	key := "books:" + idstr
	fmt.Println(key)
	data, err := RedisConn.Get(c, key).Bytes()
	if err == redis.Nil {
		// 如果Redis中没有缓存数据，则查询MySQL数据库，并将结果压缩存放到Redis中
		var books []Book
		sql := "select id,name,count,img_url from books use index(search) where id>=? limit ? order by id asc" //分页查询
		DB.Raw(sql, id, size).Find(&books)
		//压缩
		Ybyte := Yasuo(books)
		data = Ybyte
		//防止击穿这种方法可以用来避免大量的键同时过期，从而减轻 Redis 的负载压力。同时，随机过期时间也可以使键的过期时间更加均匀，
		//避免数据集中存储在某个时间段内过期，造成 Redis 的短时间内负载过高。
		num := rand.Intn(3) + 3 //3-5之间的随机数
		err = RedisConn.Set(c, key, Ybyte, time.Duration(num)*time.Second).Err()
		if err != nil {
			c.AbortWithError(500, err)
			return nil
		}
	}
	//解压
	return Jieya(data)
}

// SearchCategory 查询分类
func SearchCategory() []Classification {
	var clas []Classification
	sql := "select * from classifications"
	err := DB.Raw(sql).Find(&clas).Error
	if err != nil {
		fmt.Println("err+?", err)
		return nil
	}
	return clas
}
