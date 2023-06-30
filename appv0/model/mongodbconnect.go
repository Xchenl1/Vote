package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Client *mongo.Client

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// Connect to MongoDB mongo.Connect() 函数创建了一个客户端实例 传递context参数
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection 检查是否连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	Client = client
	usersCollection := client.Database("lol-lpl-data").Collection("match_basic")
	fmt.Println(usersCollection)
	//fmt.Println("Connected to MongoDB!")
	//创建一个过滤器  查数据 bson.A: map[string]int  bson.D: map[string]interface{}
	//查询id>=10163的数据
	//filter := bson.M{"matchid": bson.M{"$gt": 0}}
	//cur, err1 := usersCollection.Find(context.Background(), filter)
	////随手关闭
	//defer func(cur *mongo.Cursor, ctx context.Context) {
	//	err := cur.Close(ctx)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}(cur, context.Background())
	//nums := 0
	//for cur.Next(context.Background()) {
	//	var result bson.M
	//	err := cur.Decode(&result)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(result)
	//	nums++
	//}
	//if err1 != nil {
	//	fmt.Println("err1", err1)
	//}
	//fmt.Println(nums)
	//fmt.Println(result)
}
