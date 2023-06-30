package logic

import (
	"book_manage_system/appv0/model"
	"book_manage_system/appv0/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// AddUser
//
//	@Summary		用户注册
//	@Description	用户注册
//	@Tags			login
//	@Produce		json
//	@Param			userName	formData	string	true	"昵称"
//	@Param			password	formData	string	true	"密码"
//	@Param			name		formData	string	true	"姓名"
//	@Param			sex			formData	string	true	"性别"
//	@Param			phone		formData	string	true	"电话号码"
//	@Success		200,500		{object}	tools.HttpCode
//	@Router			/users [post]
func AddUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "绑定失败！",
		})
		return
	}
	k, _ := c.GetPostForm("sex")
	if k == "男" || k == "女" {
		goto begin
	} else {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "性别有误",
		})
		return
	}
begin:
	l, _ := c.GetPostForm("phone")
	l1, _ := strconv.ParseInt(l, 10, 64)
	if count := model.Weishu(l1); count != 11 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "电话号码位数不对",
		})
		return
	}
	if bool1 := model.Adduser(user.UserName, user.Pwd, user.Name, user.Sex, user.Phone); bool1 != true {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "注册失败",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "成功",
	})
	return
}

// GetUser godoc
//
//	@Summary		用户查看信息
//	@Description	用户查看自己信息
//	@Tags			user
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Success		200				{object}	tools.HttpCode
//	@Router			/users [GET]
func GetUser(c *gin.Context) {
	k, _ := c.Get("userId")
	if k.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "id错误！",
			Data:    nil,
		})
		return
	}
	user := model.Getusersd(k.(int64))
	if user == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没有找到数据",
		})
		return
	}
	//200成功
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "返回用户信息",
		Data:    user,
	})
	return
}

type Updateuser struct {
	Username    string `json:"username"  form:"username"`
	Password    string `json:"password"  form:"password"`
	Newpassword string `json:"newpassword" form:"newpassword"`
	Phone       string `json:"phone"      form:"phone"`
}

// UpdateUser godoc
//
//	@Summary		用户修改信息
//	@Description	用户修改自己的信息
//	@Tags			user
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			username		formData	string	true	"用户名"
//	@Param			password		formData	string	true	"旧密码"
//	@Param			newpassword		formData	string	true	"新密码"
//	@Param			phone			formData	string	true	"电话"
//	@Response		200,500			{object}	tools.HttpCode
//	@Router			/users/{id} [PUT]
func UpdateUser(c *gin.Context) {
	var updateuser1 Updateuser
	if err := c.ShouldBind(&updateuser1); err != nil {
		//400错误是表示请求无效，语法错误或无法满足请求
		c.JSON(http.StatusBadRequest, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "绑定更新用户信息失败！",
			Data:    nil,
		})
		return
	}
	//fmt.Println(updateuser1)
	//更新操作 写事务
	if updateuser1.Password == updateuser1.Newpassword {
		l, _ := c.GetPostForm("phone")
		l1, _ := strconv.ParseInt(l, 10, 64)
		if count := model.Weishu(l1); count != 11 {
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "电话号码位数不对",
			})
			return
		}
		K, _ := c.Get("userId")
		t := model.Updateuser(updateuser1.Username, updateuser1.Password, updateuser1.Phone, K.(int64))
		if t {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: "更新信息成功",
				Data:    nil,
			})
		} else {
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.DoErr,
				Message: "更新失败",
			})
			return
		}
	} else {
		//200成功
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "新旧密码不一致",
			Data: struct {
			}{},
		})
		return
	}
}

// GetRecords
//
//	@Summary		获取借书记录
//	@Description	获取用户借书记录
//	@Tags			book
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Param			id				path		int		true	"用户id"
//	@Success		200				{object}	tools.HttpCode{data=model.Book}
//	@Router			/users/{id}/records [GET]
func GetRecords(c *gin.Context) {
	k, _ := c.Get("userId")
	//k1, _ := strconv.ParseInt(k.(string), 10, 64)
	if k.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户",
			Data:    nil,
		})
		return
	}
	id := c.Param("id")
	id1, _ := strconv.ParseInt(id, 10, 64)
	record := model.Findrecord(id1)
	if record == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到记录",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "",
		Data:    record,
	})
	return
}

// GetStatusRecords
//
//	@Summary		获取还书记录
//	@Description	获取用户还书记录
//	@Tags			book
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			status			path		string	true	"图书状态"
//	@Param			id				path		string	true	"用户id"
//	@Success		200				{object}	tools.HttpCode{data=model.Book}
//	@Router			/users/{id}/records/{status} [GET]
func GetStatusRecords(c *gin.Context) {
	k := c.Param("id")
	l := c.Param("status")
	k1, _ := strconv.ParseInt(k, 10, 64)
	l1, _ := strconv.ParseInt(l, 10, 64)
	if k1 <= 0 && l1 <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户",
			Data:    nil,
		})
		return
	}
	var record []model.Record
	sql := "select * from records where user_id =? and status=?"
	err := model.DB.Raw(sql, k1, l1).Find(&record).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到记录",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "",
		Data:    record,
	})
	return
}

// BorrowBook
//
//	@Summary		借书
//	@Description	借书
//	@Tags			book
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			bookid			formData	string	true	"书籍id"
//	@Param			time			formData	string	true	"借书截至时间"	format("2006-01-02 15:04:05")
//	@Success		200,500			{object}	tools.HttpCode
//	@Router			/users/records/{bookid} [Post]
func BorrowBook(c *gin.Context) {
	k, _ := c.GetPostForm("bookid")
	q, _ := c.Get("userId")
	k1, _ := strconv.ParseInt(k, 10, 64)
	time1, _ := c.GetPostForm("time")
	if k1 <= 0 || q.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到记录",
			Data:    nil,
		})
		return
	}
	var user model.User
	//forupdate
	sql2 := "select * from users where id=? for update"
	err := model.DB.Raw(sql2, q.(int64)).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没找到用户",
		})
		return
	}
	if user.Id == 0 || user.Status == 1 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "不能借书",
		})
		return
	}
	sql1 := "select * from books where id=?"
	var book model.Book
	err1 := model.DB.Raw(sql1, k1).Find(&book).Error
	if err1 != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到书籍",
		})
		return
	}
	if book.Id == 0 || book.Count == 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "书籍没找到或者已被借完",
		})
		return
	}
	//事务
	tx := model.DB.Begin()
	sql := "insert into records(user_id,book_id,status,start_time,over_time) values (?,?,?,?,?)"
	err3 := tx.Exec(sql, q.(int64), k1, 0, time.Now().Format("2006-01-02 15:04:05"), time1).Error
	if err3 != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "借书失败",
		})
		return
	}
	tx.Commit()
	//改变该书籍的数量
	tx1 := model.DB.Begin()
	book.Count -= 1
	sql3 := "update books set count=? where id=?"
	err2 := tx1.Exec(sql3, book.Count, book.Id).Error
	if err2 != nil {
		tx1.Rollback()
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "更新数量出现异常！",
		})
		return
	}
	tx1.Commit()
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "借书成功",
	})
	return
}

// ReturnBook
//
//	@Summary		还书
//	@Description	还书
//	@Tags			book
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			bookid			formData	string	true	"书籍id"
//	@Success		200,500			{object}	tools.HttpCode
//	@Router			/users/records/{bookid} [PUT]
func ReturnBook(c *gin.Context) {
	k, _ := c.GetPostForm("bookid")
	k1, _ := strconv.ParseInt(k, 10, 64)
	userid, _ := c.Get("userId")
	if k1 <= 0 || userid.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找用户信息",
		})
		return
	}
	var record model.Record
	sql := "select * from records where user_id=? and book_id=? "
	err := model.DB.Raw(sql, userid.(int64), k1).Find(&record).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户信息",
		})
		return
	}
	tx := model.DB.Begin()
	sql1 := "update records set return_time=?,status=?"
	//1表示已归还
	err1 := tx.Exec(sql1, time.Now().Format("2006-01-02 15:04:05"), 1).Error
	if err1 != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "更新失败",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "还书成功！",
	})
	tx.Commit()
	return
}

// GetBook godoc
//
//	@Summary		获取图书信息
//	@Description	获取一个图书的信息
//	@Tags			book
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			id				path		string	true	"书籍id"
//	@Success		200				{object}	tools.HttpCode{data=model.Book}
//	@Router			/users/books/{id} [GET]
func GetBook(c *gin.Context) {
	l := c.Param("id")
	l1, _ := strconv.ParseInt(l, 10, 64)
	if l1 <= 0 {
		//404
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取图书信息错误！",
			Data:    nil,
		})
		return
	}
	book := model.FindbookId(l1)
	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "信息有误",
		})
		return
	}
	if book.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "不存该书",
		})
		return
	}
	book.Img_url = "/view/" + book.Img_url
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "找到书籍",
		Data:    book,
	})
	return
}

// GetCategoryBooks
//
//	@Summary		根据分类id查书
//	@Description	根据种类id查找书籍
//	@Tags			user
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Param			id				path		string	true	"书籍种类id"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/users/categories/{id} [GET]
func GetCategoryBooks(c *gin.Context) {
	k := c.Param("id")
	fmt.Println(k)
	k1, _ := strconv.ParseInt(k, 10, 64)
	Classification := model.ChafenleiId(k1)
	if Classification == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查询语句有误",
		})
		return
	}
	if Classification.Id == 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到数据",
		})
		return
	}
	var book *[]model.Book
	book = model.FindbookcateId(Classification.Id)
	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "有错误",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "查询成功！",
		Data:    book,
	})
}

// Getbooks
//
//	@Summary		用户浏览书籍
//	@Description	用户浏览书籍
//	@Tags			user
//	@Produce		json
//	@Param			id				query		string	false	"起始id"
//	@Param			size			query		string	false	"每页书籍数量"
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Response		200,500			{object}	tools.HttpCode
//	@Router			/users/books/page [get]
func Getbooks(c *gin.Context) {
	idstr := c.DefaultQuery("id", "1")
	sizeStr := c.DefaultQuery("size", "100")
	id, err := strconv.Atoi(idstr)
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
	book := model.SelectBook(c, idstr, id, size)
	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未查到数据",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "成功！",
		Data:    book,
	})
	return
}
