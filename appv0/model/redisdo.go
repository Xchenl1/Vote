package model

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"strconv"
	"time"
)

var RedisConn *redis.Client

func init() {
	// 创建 Redis 客户端连接
	RedisConn = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "qq74263827", // Redis 未设置密码时为空
		Password: "", // Redis 未设置密码时为空
		DB:       1,  // 使用默认数据库
	})
	// 测试连接是否成功
	_, err := RedisConn.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v", err)
		return
	}
	//预热缓存
	//preWarm()
	fmt.Println("Connected to Redis")
}

//gocatch 实现  解决数据不一致的问题 1.多级catch 2. 缓存一致性

// Yasuo 压缩book切片
func Yasuo(book []Book) []byte {
	// 创建 gzip 编码器 如何压缩的？
	//字节缓冲区
	var buf bytes.Buffer
	//创建一个gzip.writer对象
	gzipWriter := gzip.NewWriter(&buf)
	// 将 books 切片转换为 JSON 数据
	jsonBytes, err := json.Marshal(&book)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) :%+v\n", err.Error())
	}
	//写入缓冲区
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write :%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close :%+v\n", err.Error())
	}
	//最终得到的压缩后的数据存储在buf中
	return buf.Bytes()
}

// Jieya
func Jieya(key []byte) []Book {
	// 创建 gzip 解码器
	gzipReader, err := gzip.NewReader(bytes.NewReader(key))
	if err != nil {
		fmt.Printf("vgzip.NewReader 时出现错误！err:%+v\n", err.Error())
		return nil
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			fmt.Printf("gzipReader.Close() 时出现错误！err:%+v\n", err.Error())
			return
		}
	}(gzipReader)
	//
	books := make([]Book, 0)
	// 分批读取解压缩后的 JSON 数据
	decoder := json.NewDecoder(gzipReader)
	batchSize := 100 // 每次读取的批次大小为 100 条记录
	for {
		var batch []Book
		err := decoder.Decode(&batch)
		if err == io.EOF { // 已经读取完数据
			break
		} else if err != nil {
			fmt.Printf("ecoder.Decode(&batch) 时出现错误！err:%+v\n", err.Error())
			return nil
		}
		books = append(books, batch...)
		if len(books) >= batchSize { // 达到批次大小，返回结果
			break
		}
	}
	for i := 0; i < len(books); i++ {
		books[i].Img_url = "/view/" + books[i].Img_url
	}
	return books
}

// Loader 缓存预热
func Loader(id int, size int) []Book {
	var books []Book
	sql := "select * from books where id>=? limit ?" //分页查询
	DB.Raw(sql, id, size).Find(&books)
	return books
}

// Handler 处理热点数据没3秒更新前三页数据
func Handler(id int, books []Book) {
	c := context.Background()
	v := Yasuo(books)
	id1 := strconv.Itoa(id)
	key := "books:" + id1
	err := RedisConn.Set(c, key, v, 3*time.Second).Err()
	if err != nil {
		fmt.Printf("err", err)
	}
}
