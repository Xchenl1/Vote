package model

import (
	"fmt"
)

// Userduibi 对比密码
func Userduibi(name string, pwd string) *User {
	user1 := User{}
	sql := "select * from users where user_name = ? and pwd = ? "
	err := DB.Raw(sql, name, pwd).Scan(&user1).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &user1
}

// Adduser 添加用户
func Adduser(user_name, pwd, name, sex, phone string) bool {
	tx := DB.Begin()
	sql := "insert into users(user_name,pwd,name,sex,phone) values(?,?,?,?,?)"
	err := tx.Exec(sql, user_name, pwd, name, sex, phone).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("err", err)
		return false
	}
	tx.Commit()
	return true
}

// Weishu 计算位数
func Weishu(n int64) int64 {
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	return int64(count)
}

// ChafenleiId 根据分类id查分类表
func ChafenleiId(id int64) *Classification {
	var Classification1 Classification
	sql := "select * from classifications where id = ? "
	err := DB.Raw(sql, id).Find(&Classification1).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Classification1
}

// FindbookId 根据分类id查书
func FindbookcateId(id int64) *[]Book {
	var book []Book
	sql1 := "select * from books where classification_id=? "
	err1 := DB.Raw(sql1, id).Find(&book).Error
	if err1 != nil {
		fmt.Println("err1:", err1)
		return nil
	}
	return &book
}

// Updateuser 用户修改个人信息
func Updateuser(username string, passwd string, phone string, id int64) bool {
	tx := DB.Begin()
	sql := "update users set user_name=?,pwd=?,phone=? where id =?"
	var user2 User
	//不知道为什么用raw就不行 exec就可以  raw通常是select 防止sql注入 增删改通常是exec
	err := tx.Exec(sql, username, passwd, phone, id).Find(&user2).Error
	if err != nil {
		fmt.Println("err", err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// 获取个人信息
func Getusersd(id int64) *User {
	var user User
	sql := "select * from users where id=?"
	err := DB.Raw(sql, id).Find(&user).Error
	if err != nil {
		fmt.Println("err", err)
		return nil
	}
	return &user
}

// 根据id查找借书记录
func Findrecord(id int64) *[]Record {
	var record []Record
	sql := "select * from records where id=?"
	err := DB.Raw(sql, id).Find(&record).Error
	if err != nil {
		fmt.Println("err", err)
		return nil
	}
	return &record
}
