package logic

import (
	"book_manage_system/appv0/model"
	"book_manage_system/appv0/tools"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Name string `json:"name" form:"name" binding:"required"`
	Pwd  string `json:"pwd" form:"pwd"  binding:"required"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserLogin
//
//	@Summary		用户登录
//	@Description	执行用户登录操作
//	@Tags			login
//	@Accept			multipart/form-data
//	@Param			name	formData	string	true	"用户名"
//	@Param			pwd		formData	string	true	"密码"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/userLogin [POST]
func UserLogin(c *gin.Context) {
	//普通用户登录
	User1 := User{}
	//获取前端的登录数据
	if err := c.ShouldBind(&User1); err != nil {
		c.JSON(http.StatusBadRequest, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "数据绑定失败",
			Data: struct {
			}{},
		})
		return
	}
	//var user2 User
	data, _ := model.RedisConn.Get(c, User1.Name).Result()
	var user3 User
	_ = json.Unmarshal([]byte(data), &user3)
	fmt.Println(user3)
	//没有找到就去数据库中查找
	if user3.Name == "" {
		user2 := model.Userduibi(User1.Name, User1.Pwd)
		if user2.Id > 0 {
			//直接将id塞进去 然后就过滤 防止网站攻击
			c.SetCookie("id", strconv.FormatInt(user2.Id, 10), 3600, "/", "", false, true)
			a, r, err := tools.Token.GetToken(user2.Id, user2.UserName)
			if err != nil {
				c.JSON(http.StatusUnauthorized, tools.HttpCode{
					Code:    tools.UserInfoErr,
					Message: "Token生成失败！",
					Data: struct {
					}{},
				})
				return
			}
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: "发放Token",
				Data: Token{
					RefreshToken: r,
					AccessToken:  a,
				},
			})
			//必须转换为json类型才能存储
			k, err := json.Marshal(User1)
			if err != nil {
				fmt.Println("反序列化失败！")
			}
			model.RedisConn.Set(c, User1.Name, k, 1*time.Hour)
			return
		}
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "数据验证错误",
			Data: struct {
			}{},
		})
		return
		//在redis中找到数据之后可以直接对比
	} else {
		if user3.Pwd == User1.Pwd {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: "成功登录",
			})
			return
		} else {
			c.JSON(http.StatusNotFound, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "密码输入错误",
			})
			return
		}
	}
}

// LibrarianLogin
//
//	@Summary		管理员登录
//	@Description	执行管理员登录操作
//	@Tags			login
//	@Accept			multipart/form-data
//	@Param			name	formData	string	true	"用户名"in:formData
//	@Param			pwd		formData	string	true	"密码"in:formData
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/adminLogin [POST]
func LibrarianLogin(c *gin.Context) {
	//复用用户登录的代码
	User1 := User{}
	if err := c.ShouldBind(&User1); err != nil {
		c.JSON(http.StatusBadRequest, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "数据绑定失败",
			Data: struct {
			}{},
		})
		return
	}
	adminstr := model.Adminlogin(User1.Name, User1.Pwd)
	if adminstr.Id == 0 {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未找到管理员",
		})
		return
	}
	if adminstr == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查找出错",
		})
		return
	}
	err1 := tools.SetSession(c, adminstr.Name, adminstr.Id)
	if err1 != nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "获取session失败",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "登陆成功",
	})
	return
}
