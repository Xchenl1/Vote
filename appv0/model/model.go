package model

import "time"

// Book 图书表
type Book struct {
	Id               int64  `json:"id"`                                    //主键id
	BN               string `json:"bn"gorm:"type:varchar(20)"`             //图书编号
	Name             string `json:"name"gorm:"type:varchar(200)"`          //图书名
	Description      string `json:"description"gorm:"type:varchar(15000)"` //图书描述
	Count            int    `json:"count"`                                 //数量 默认为1
	ClassificationId int64  `json:"classificationId"`                      //分类id
	Img_url          string `json:"img_Url" gorm:"varchar(200)"`           //图片地址
}

type BookInfo struct {
	ID                 uint      `gorm:"column:id"`
	BookName           string    `gorm:"column:book_name"`
	Author             string    `gorm:"column:author"`
	PublishingHouse    string    `gorm:"column:publishing_house"`
	Translator         string    `gorm:"column:translator"`
	PublishDate        time.Time `gorm:"column:publish_date"`
	Pages              int       `gorm:"column:pages"`
	ISBN               string    `gorm:"column:ISBN"`
	Price              float64   `gorm:"column:price"`
	BriefIntroduction  string    `gorm:"column:brief_introduction"`
	AuthorIntroduction string    `gorm:"column:author_introduction"`
	ImgURL             string    `gorm:"column:img_url"`
	DelFlg             int       `gorm:"column:del_flg"`
}

// classification 分类表
type Classification struct {
	Id       int64  `json:"id"`       //主键id
	Category string `json:"category"` //分类
	Book     []Book `json:"book"`     //存放这类书
}

// User 用户表
type User struct {
	Id       int64  `json:"id" form:"id"`                                      //主键id
	UserName string `json:"userName" form:"userName" gorm:"type:varchar(100)"` //表是该字段的类型是varchar 用户名
	Pwd      string `json:"password" form:"password" gorm:"type:varchar(100)"` //密码
	Name     string `json:"name" form:"name" gorm:"type:varchar(100)"`         //姓名
	Sex      string `json:"sex" form:"sex" gorm:"type:varchar(100)"`           //性别
	Phone    string `json:"phone" form:"phone" gorm:"type:varchar(100)"`       //电话号码
	Status   int    `json:"status" form:"status"`                              //0:正常  1:封禁 状态
}

// Librarian 图书管理员
type Librarian struct {
	Id       int64  //主键id
	UserName string `gorm:"type:varchar(100)"` //管理员昵称
	Password string `gorm:"type:varchar(100)"` //密码
	Name     string `gorm:"type:varchar(100)"` //真实姓名
	Sex      string `gorm:"type:varchar(100)"` //性别
	Phone    string `gorm:"type:varchar(100)"` //电话号码
}

// Record 借书表
type Record struct {
	Id         int64     //主键id
	UserId     int64     //用户id
	BookId     int64     //书id
	Status     int       //已归还1 未归还0
	StartTime  time.Time //借书时间
	OverTime   time.Time //
	ReturnTime time.Time //还书时间
}

// 定时器发送信息
type Sendmessage struct {
	Id        int64     `json:"id"`
	Userid    int64     `json:"userid"`
	Bookid    int64     `json:"bookid"`
	Message   string    `json:"message"`
	Starttime time.Time `json:"starttime"`
}
