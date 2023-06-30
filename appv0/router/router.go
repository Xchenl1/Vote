package router

import (
	"book_manage_system/appv0/logic"
	"book_manage_system/appv0/tools"
	_ "book_manage_system/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func New() *gin.Engine {
	//定时器
	go tools.TimerMiddleware()
	//加载热点信息 定时器每三秒加载前三页书籍信息
	//go tools.Cacheheating()

	r := gin.Default()
	r.Static("/view", "./view")
	//登录模块
	{
		//用户登录
		r.POST("/userLogin", logic.UserLogin)
		//用户注册
		r.POST("/users", logic.AddUser)
		//获取验证码
		r.GET("/sendcode", tools.Sendcode)
		//验证码登录
		r.POST("/codelogin", tools.Codelogin)
		//管理员登录
		r.POST("/adminLogin", logic.LibrarianLogin)
		//游客可以浏览书籍
		r.GET("/books", logic.SearchBook)
		//游客可以浏览分类
		r.GET("/categories", logic.SearchCategory)
	}

	//用户模块
	user := r.Group("users")
	//使用Token中间件 防止链接被多次点击
	user.Use(tools.UserAuthCheck(), tools.LimitedFlow(2, 2*time.Second))
	{
		//获取用户个人信息 实现
		user.GET("", logic.GetUser)
		//修改个人信息 实现
		user.PUT("/:id", logic.UpdateUser)
		//获取个人借书记录 实现
		user.GET("/:id/records", logic.GetRecords)
		//获取用户个人还书记录 实现
		user.GET("/:id/records/:status", logic.GetStatusRecords)
		//用户借书 实现
		user.POST("/records/:bookId", logic.BorrowBook)
		//用户还书 实现
		user.PUT("/records/:bookId", logic.ReturnBook)

		book := user.Group("/books")
		{
			//查看书的具体信息 实现
			book.GET("/:id", logic.GetBook)
			//获取所有书籍信息
			book.GET("/page", logic.Getbooks)
		}
		category := user.Group("/categories")
		{
			//根据分类查书 实现
			category.GET("/:id", logic.GetCategoryBooks)
			//根据分类名查询图书 实现
			//category.GET("", logic.Getcategorybookname)
		}
	}

	//管理员模块
	admin := r.Group("administer")
	//管理员中间件
	admin.Use(tools.AdminAuthCheck(), tools.LimitedFlow(2, 2*time.Second))
	{
		user1 := admin.Group("/users")
		{
			//获取用户信息 实现
			user1.GET("", logic.SearchUser)
			//修改用户信息 实现
			user1.PUT("/:id", logic.UpdateUserByAdmin)
			//删除用户信息 实现
			user1.DELETE("/:id", logic.DeleteUser)
			//查看记录表 实现
			user1.GET("/:id/records/:status", logic.GetUserBook)
		}
		//书的所有资源
		book1 := admin.Group("/books")
		{
			//获取书的详细记录 实现
			book1.GET("/:id", logic.GetBook1)
			//添加书 实现
			book1.POST("", logic.AddBook)
			//更新书 实现
			book1.PUT("/:id", logic.UpdateBook)
			//删除书 实现
			book1.DELETE("/:id", logic.DeleteBook)
			//获取所有书籍
			book1.GET("/page", logic.Getbooks1)
		}
		category := admin.Group("/categories")
		{
			//查看图书分类 实现
			category.GET("", logic.GetCategory)
			//添加图书分类 实现
			category.POST("", logic.AddCategory)
			//修改图书分类 实现
			category.PUT("/:id", logic.UpdateCategory)
			//删除图书分类 实现
			category.DELETE("/:id", logic.DeleteCategory)
		}
		record := admin.Group("/records")
		{
			//所有借书还书记录
			record.GET("", logic.GetRecords1)
			//所有归还或者未归还的记录
			record.GET("/:status", logic.GetRecords2)
		}
	}
	//swagger测试
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
