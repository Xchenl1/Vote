package logic

import (
	"book_manage_system/appv0/model"
	"book_manage_system/appv0/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SearchUser
//
//	@Summary		管理员获取用户信息
//	@Description	管理员获取用户信息
//	@Tags			adminstr
//	@Produce		json
//	@Response		200	{object}	tools.HttpCode
//	@Router			/administer/users [GET]
func SearchUser(c *gin.Context) {
	sql := "select * from users"
	var user []model.User
	err := model.DB.Raw(sql).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户信息",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "成功！",
		Data:    user,
	})
	return
}

// UpdateUserByAdmin
//
//	@Summary		管理员修改用户信息
//	@Description	管理员修改用户信息
//	@Tags			adminstr
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		string	true	"用户id"
//	@Param			username	formData	string	true	"用户名"
//	@Param			password	formData	string	true	"旧密码"
//	@Param			newpassword	formData	string	true	"新密码"
//	@Param			phone		formData	string	true	"电话"
//	@Response		200,500		{object}	tools.HttpCode
//	@Router			/administer/users/{id} [PUT]
func UpdateUserByAdmin(c *gin.Context) {
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
	K := c.Param("id")
	k1, _ := strconv.ParseInt(K, 10, 64)
	var user model.User
	sql := "select * from users where id=?"
	err := model.DB.Raw(sql, k1).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查询语句有误",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "用户不存在",
		})
		return
	}
	if updateuser1.Password == updateuser1.Newpassword {
		tx := model.DB.Begin()
		sql := "update users set user_name=?,pwd=?,phone=? where id =?"
		//fmt.Println(K)
		var user2 model.User
		//不知道为什么用raw就不行 exec就可以  raw通常是select 防止sql注入 增删改通常是exec
		err := tx.Exec(sql, updateuser1.Username, updateuser1.Password, updateuser1.Phone, k1).Find(&user2).Error
		if err != nil {
			//404
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.DoErr,
				Message: "更新用户信息失败",
				Data:    nil,
			})
			tx.Rollback()
			return
		}
		//200成功
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "更新信息成功",
			Data:    nil,
		})
		tx.Commit()
		return
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

// DeleteUser
//
//	@Summary		管理员删除用户
//	@Description	管理员删除用户信息
//	@Tags			adminstr
//	@Produce		json
//	@Param			id		path		string	true	"用户id"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/users/{id} [DELETE]
func DeleteUser(c *gin.Context) {
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	if k1 <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户",
		})
		return
	}
	sql := "delete from users where id=?"
	tx := model.DB.Begin()
	err := tx.Exec(sql, k1).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "删除用户失败",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "删除用户成功",
	})
	tx.Commit()
	return
}

// GetUserBook
//
//	@Summary		管理员查看借书信息
//	@Description	管理员信息
//	@Tags			adminstr
//	@Produce		json
//	@Param			id		path		string	true	"用户id"
//	@Param			status	path		string	true	"图书状态"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/users/{id}/records/{status} [Get]
func GetUserBook(c *gin.Context) {
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	l := c.Param("status")
	l1, _ := strconv.ParseInt(l, 10, 64)
	if k1 <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户信息",
		})
		return
	}
	if l1 == 1 || l1 == 0 {
		var record []model.Record
		sql := "select * from records where user_id=? and status=?"
		err := model.DB.Raw(sql, k1, l1).Find(&record).Error
		if err != nil {
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "未找到数据",
			})
			return
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "成功",
			Data:    record,
		})
		return
	} else {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "错误",
		})
		return
	}
}

// GetBook1  godoc
//
//	@Summary		获取图书信息
//	@Description	获取一个图书的信息
//	@Tags			adminstr
//	@Produce		json
//	@Param			id	path		string	true	"书籍id"
//	@Success		200	{object}	tools.HttpCode
//	@Router			/administer/books/{id} [GET]
func GetBook1(c *gin.Context) {
	l := c.Param("id")
	l1, _ := strconv.ParseInt(l, 10, 64)
	book := model.FindbookId(l1)
	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查询语句有误",
		})
		return
	}
	if book.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没找都书籍！",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "找到书籍",
		Data:    book,
	})
	return
}

type TBook struct {
	Bn          string `json:"bn" form:"bn"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Count       int    `json:"count" form:"count"`
	Category    int    `json:"category" form:"category"`
}

// AddBook
//
//	@Summary		添加图书
//	@Description	添加图书的信息
//	@Tags			adminstr
//	@Produce		json
//	@Param			bn			formData	string	true	"图书编号"
//	@Param			name		formData	string	true	"书名"
//	@Param			description	formData	string	true	"描述"
//	@Param			count		formData	int		true	"数量"
//	@Param			category	formData	int		true	"类别"
//	@Response		200,500		{object}	tools.HttpCode
//	@Router			/administer/books [post]
func AddBook(c *gin.Context) {
	var book TBook
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定数据失败",
		})
		return
	}
	cate, _ := c.GetPostForm("category")
	cate1, _ := strconv.ParseInt(cate, 10, 64)
	n := model.FindbookId(cate1)
	if n == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "有误",
		})
		return
	}
	if n.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到分类",
		})
		return
	}
	s := model.Addbook(book.Bn, book.Name, book.Description, book.Count, book.Category)
	if s == false {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "成功",
	})
	return
}

// UpdateBook
//
//	@Summary		修改图书信息
//	@Description	修改图书信息
//	@Tags			adminstr
//	@Produce		json
//	@Param			id			path		string	true	"书籍id"
//	@Param			bn			formData	string	true	"图书编号"
//	@Param			name		formData	string	true	"书名"
//	@Param			description	formData	string	true	"描述"
//	@Param			count		formData	int		true	"数量"
//	@Param			category	formData	int		true	"类别"
//	@Response		200,500		{object}	tools.HttpCode
//	@Router			/administer/books/{id} [put]
func UpdateBook(c *gin.Context) {
	var book TBook
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定数据失败",
		})
		return
	}
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	book1 := model.FindbookId(k1)
	if book1 == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到数据",
		})
		return
	}
	if book1.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没有该书",
		})
		return
	}

	//判断种类是否存在
	zl, _ := c.GetPostForm("category")
	zl1, _ := strconv.ParseInt(zl, 10, 64)
	category := model.ChafenleiId(zl1)
	if category.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "种类输入错误",
		})
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到结果",
		})
		return
	}
	//看种类是否有错误
	sl, _ := c.GetPostForm("count")
	sl1, _ := strconv.ParseInt(sl, 10, 64)
	if sl1 < 0 || sl1 >= 1000 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "数量有误",
		})
		return
	}
	l := model.Updatebook(book.Bn, book.Name, book.Description, book.Count, k1)
	if l {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "更新书籍成功",
		})
		return
	} else {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "更新书籍出错",
		})
		return
	}
}

// DeleteBook
//
//	@Summary		删除图书
//	@Description	删除图书
//	@Tags			adminstr
//	@Produce		json
//	@Param			id		path		string	true	"书籍id"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	book := model.FindbookId(k1)
	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查询语句有误",
		})
		return
	}
	if book.Id == 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "书籍没找到",
		})
		return
	}
	bol := model.Deletebook(k1)
	if bol {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "删除成功！",
		})
		return
	} else {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "删除失败",
		})
		return
	}
}

// GetCategory
//
//	@Summary		获取分类
//	@Description	获取分类
//	@Tags			adminstr
//	@Produce		json
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/categories [get]
func GetCategory(c *gin.Context) {
	var book []model.Classification
	sql := "select * from classifications "
	err := model.DB.Raw(sql).Find(&book).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没有分类",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "返回分类列表",
		Data:    book,
	})
	return
}

// AddCategory
//
//	@Summary		添加分类
//	@Description	添加分类
//	@Tags			adminstr
//	@Produce		json
//	@Param			category	formData	string	true	"分类名"
//	@Response		200,500		{object}	tools.HttpCode
//	@Router			/administer/categories [post]
func AddCategory(c *gin.Context) {
	k, _ := c.GetPostForm("category")
	var classification model.Classification
	sql := "select *  from  classifications where category=?"
	err := model.DB.Raw(sql, k).Find(&classification).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没有找到数据",
		})
		return
	}
	//不能重复
	if classification.Category == k {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "分类名重复",
		})
		return
	} else {
		tx := model.DB.Begin()
		sql1 := "insert into classifications(category) values (?)"
		err1 := tx.Exec(sql1, k).Error
		if err1 != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "添加失败",
			})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "插入成功",
		})
		return
	}
}

// UpdateCategory
//
//	@Summary		修改分类
//	@Description	修改分类
//	@Tags			adminstr
//	@Produce		json
//	@Param			id			path		int64	true	"分类id"
//	@Param			category	formData	string	true	"分类名"
//	@Response		200,500		{object}	tools.HttpCode
//	@Router			/administer/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	q, _ := c.GetPostForm("category")
	var classification model.Classification
	sql := "select *  from  classifications where category=?"
	err := model.DB.Raw(sql, q).Find(&classification).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "查询错误",
		})
		return
	}
	if classification.Id > 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "分类名重复",
		})
		return
	} else {
		tx := model.DB.Begin()
		sql1 := "update classifications set category=? where id=?"
		err1 := tx.Exec(sql1, q, k1).Error
		if err1 != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "更新分类失败",
			})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "更新成功",
		})
		return
	}

}

// DeleteCategory
//
//	@Summary		删除分类
//	@Description	删除分类
//	@Tags			adminstr
//	@Produce		json
//	@Param			id		path		int64	true	"分类id"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	k := c.Param("id")
	k1, _ := strconv.ParseInt(k, 10, 64)
	if k1 <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "数据错误",
		})
		return
	}
	var clas model.Classification
	sql1 := "select * from classifications where id=?"
	err1 := model.DB.Raw(sql1, k1).Find(&clas).Error
	if err1 != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查找失败",
		})
		return
	}
	if clas.Id <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "分类不存在！",
		})
		return
	}
	tx := model.DB.Begin()
	sql := "delete from classifications where id=?"
	err := tx.Exec(sql, k1).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "删除失败",
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "删除成功！",
	})
	return
}

// GetRecords1
//
//	@Summary		查看所有记录
//	@Description	查看所有记录
//	@Tags			adminstr
//	@Produce		json
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/records [get]
func GetRecords1(c *gin.Context) {
	k, _ := c.Get("userId")
	if k.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "没有找到用户信息",
		})
		return
	}
	var record []model.Record
	sql := "select * from records"
	err := model.DB.Raw(sql).Find(&record).Error
	if err != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到记录信息",
		})
		return
	}
	if len(record) == 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "没有记录",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "找到数据,已返回",
		Data:    record,
	})
	return
}

// GetRecords2
//
//	@Summary		查看所有归还或未归还记录
//	@Description	查看所有归还或未归还记录
//	@Tags			adminstr
//	@Produce		json
//	@Param			status	path		string	true	"状态"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/records/{status} [get]
func GetRecords2(c *gin.Context) {
	k := c.Param("status")
	k1, _ := strconv.ParseInt(k, 10, 64)
	l, _ := c.Get("userId") //any类型 需要断言
	if l.(int64) <= 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到用户",
		})
		return
	}
	if k1 == 0 || k1 == 1 {
		var record []model.Record
		sql := "select * from records where status=?"
		err := model.DB.Raw(sql, k1).Find(&record).Error
		if err != nil {
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "未找到信息",
			})
			return
		}
		if len(record) == 0 {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: "数据为空",
				Data:    record,
			})
			return
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "返回数据",
			Data:    record,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "数据参数不对",
		})
		return
	}
}

// Getbooks1
//
//	@Summary		管理员浏览书籍
//	@Description	管理员浏览书籍
//	@Tags			adminstr
//	@Produce		json
//	@Param			id		query		int	false	"起始页码"
//	@Param			size	query		int	false	"每页记录数"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/administer/books/page [get]
func Getbooks1(c *gin.Context) {
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
	if book != nil {
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
	} else {
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
}
