package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type Book struct {
	ID     int    `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	Author string `gorm:"not null"`
}

func main() {
	// 创建Gin实例
	r := gin.Default()

	// 连接MySQL数据库
	dsn := "root:password@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 连接Redis数据库
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 定义路由
	r.GET("/books", func(c *gin.Context) {
		// 解析分页参数
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "100")
		//转换为int
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid page number"})
			return
		}
		//转化为int
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid page size"})
			return
		}
		// 从Redis中获取分页数据
		key := "books:" + pageStr
		fmt.Println(key)
		data, err := rdb.Get(c, key).Bytes()
		if err == redis.Nil {
			// 如果Redis中没有缓存数据，则查询MySQL数据库，并将结果压缩存放到Redis中
			var books []Book
			//分页查询
			db.Offset((page - 1) * size).Limit(size).Find(&books)
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, err = gz.Write([]byte(booksString(books)))
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			err = gz.Close()
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			err = rdb.Set(c, key, buf.Bytes(), 30*time.Minute).Err()
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

		} else if err != nil {
			c.AbortWithError(500, err)
			return
		}
		// 返回结果
		c.Data(200, "application/json", data)
	})

	// 异步更新数据
	go func() {
		for {
			c := context.Background() // 创建上下文
			// 每30秒更新一次数据
			time.Sleep(30 * time.Second)

			// 从MySQL中获取最新数据
			var books []Book
			db.Find(&books)

			// 将最新数据压缩存放到Redis中
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, err = gz.Write([]byte(booksString(books)))
			if err != nil {
				log.Fatal(err)
			}
			err = gz.Close()
			if err != nil {
				log.Fatal(err)
			}

			err = rdb.Set(c, "books:latest", buf.Bytes(), 0).Err()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// 启动Web服务
	r.Run(":8080")
}

func booksString(books []Book) string {
	//var buf bytes.Buffer
	//buf.WriteString("[")
	//for i, book := range books {
	//	if i > 0 {
	//		buf.WriteString(",")
	//	}
	//	buf.WriteString(`{"id":`)
	//	buf.WriteString(strconv.Itoa(book.ID))
	//	buf.WriteString(`,"name":"`)
	//	buf.WriteString(book.Name)
	//	buf.WriteString(`","author":"`)
	//	buf.WriteString(book.Author)
	//	buf.WriteString(`"}`)
	//}
	//buf.WriteString("]")
	//ctx := context.Background() // 创建上下文
	// 创建 gzip 编码器
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	// 将 books 切片转换为 JSON 数据，并压缩
	jsonBytes, err := json.Marshal(&books)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) :%+v\n", err.Error())
	}
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write :%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close :%+v\n", err.Error())
	}
	return buf.String()
}
