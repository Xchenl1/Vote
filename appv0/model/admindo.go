package model

import (
	"fmt"
)

// 根据主键id查看图书
func FindbookId(id int64) *Book {
	var book Book
	sql := "select * from books where id=? "
	if err := DB.Raw(sql, id).Find(&book).Error; err != nil {
		fmt.Println("err", err)
		return nil
	}
	return &book
}

// 添加书
func Addbook(bn string, name string, description string, count int, category int) bool {
	tx := DB.Begin()
	sql := "insert into books(bn,name,description,count,classification_id)values(?,?,?,?,?)"
	err := tx.Exec(sql, bn, name, description, count, category).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("err", err)
		return false
	}
	tx.Commit()
	return true
}

// 更新书籍
func Updatebook(bn string, name string, description string, count int, id int64) bool {
	tx := DB.Begin()
	sql := "update books set bn=?,name=?,description=?,count=? where id=?"
	err := tx.Exec(sql, bn, name, description, count, id).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("err", err)
		return false
	}
	tx.Commit()
	return true
}

// 删除图书
func Deletebook(id int64) bool {
	tx := DB.Begin()
	sql1 := "delete from books where id=?"
	err1 := tx.Exec(sql1, id).Error
	if err1 != nil {
		tx.Rollback()
		fmt.Println("err", err1)
		return false
	}
	tx.Commit()
	return true
}

// 管理员登录
func Adminlogin(username, passwd string) *Librarian {
	var adminstr Librarian
	sql := "select * from librarians where user_name=? and password=?"
	err := DB.Raw(sql, username, passwd).Find(&adminstr).Error
	if err != nil {
		fmt.Println("err", err)
		return nil
	}
	return &adminstr
}
