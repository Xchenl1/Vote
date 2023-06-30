package appv0

import (
	"book_manage_system/appv0/model"
	"book_manage_system/appv0/router"
	"book_manage_system/appv0/tools"
)

// Start model.DB.AutoMigrate(model.Book{})
// var bookinfo []model.BookInfo
// sql := "select * from book_info"
// err := model.DB.Raw(sql).Find(&bookinfo).Error
//
//	if err != nil {
//		fmt.Println("err", err)
//	}
//
// //fmt.Println(bookinfo, len(bookinfo))
// i := 0
// k := 100
//
//	for i < len(bookinfo) {
//		rand.Seed(time.Now().UnixNano())
//		// 生成一个0到100之间的随机整数
//		randomNumber := rand.Intn(10) + 1
//		tx := model.DB.Begin()
//		sql1 := "insert into books(bn,name,description,count,classification_id,img_url) values (?,?,?,?,?,?)"
//		err1 := tx.Exec(sql1, bookinfo[i].ISBN, bookinfo[i].BookName, bookinfo[i].BriefIntroduction, k, randomNumber, bookinfo[i].ImgURL).Error
//		if err1 != nil {
//			fmt.Println(err1)
//			panic(err1)
//		}
//		tx.Commit()
//		i++
//	}
func Start() {
	//及时关闭
	defer model.Gb()
	//连接数据库
	model.Lianjie()
	//model.DB.AutoMigrate(model.Sendmessage{})
	//添加密钥
	tools.NewToken("CLT")
	//启动服务
	k := router.New()
	//监听
	k.Run(":8080")
}
