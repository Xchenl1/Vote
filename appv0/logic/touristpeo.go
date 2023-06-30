package logic

import (
	"book_manage_system/appv0/model"
	"book_manage_system/appv0/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SearchBook
//
//	@Summary		游客浏览书籍
//	@Description	游客浏览书籍
//	@Tags			tourist
//	@Produce		json
//	@Param			id		query		string	false	"起始id"
//	@Param			size	query		string	false	"每页书籍数量"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/books [get]
func SearchBook(c *gin.Context) {
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

// SearchCategory
//
//	@Summary		游客浏览分类
//	@Description	游客浏览分类
//	@Tags			tourist
//	@Produce		json
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/categories [get]
func SearchCategory(c *gin.Context) {
	cate := model.SearchCategory()
	if cate == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "未查询数据",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "查询成功！",
		Data:    cate,
	})
	return
}
