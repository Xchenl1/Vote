package model

import (
	"fmt"
	"sync"
	"time"
)

func Sendmesg() []Record {
	var record []Record
	time := time.Now().Format("2006-01-02 15:04:05")
	sql := "select * from records where datediff(over_time,?) <=2"
	err := DB.Raw(sql, time).Find(&record).Error
	fmt.Println(time)
	fmt.Println(record)
	if err != nil {
		fmt.Println("err", err)
	}
	return record
}

// Cr 插入数据库未还书的用户
func Cr(sendmesg []Record) {
	var sum int64 = 0
	var mutex sync.Mutex
	for _, v := range sendmesg {
		//这里默认的是一个人只一本相同的书
		sql1 := "select count(*) from sendmessages where userid=? and bookid=? "
		err1 := DB.Raw(sql1, v.UserId, v.BookId).Find(&sum).Error
		fmt.Println("sum", sum)
		if err1 != nil {
			fmt.Println("err1", err1)
			break
		}
		//已有记录不用再次查询
		if sum > 0 {
			continue
		}
		//已归还
		if v.Status == 1 {
			continue
		}
		//上锁
		mutex.Lock()
		tx := DB.Begin()
		sql := "insert into sendmessages(userid,message,starttime,bookid) values (?,?,?,?)"
		time1 := time.Now().Format("2006-01-02 15:04:05")
		time4 := v.OverTime.Sub(time.Now())
		mes := "距离还书时间还有" + time4.String()
		//mes := fmt.Sprintf("距离还书时间还有%time", v.OverTime-time)
		err := tx.Exec(sql, v.UserId, mes, time1, v.BookId).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("err", err)
		}
		tx.Commit()
		//解锁
		mutex.Unlock()
	}
}
