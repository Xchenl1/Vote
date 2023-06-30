package tools

import (
	"book_manage_system/appv0/model"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

//如何检验手机号
//画发送手机验证码流程图
//如果发送不成功 怎么办 如何处理？
//限制发送验证码次数
//通过配置调整1个手机号发送验证次数
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
//LTAI5t95WtavqP8qLBAmnC3N
//NiNEdINBsXcZV3Moj3wl3xkxKJPblF

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// Sendcode
//
//	@Tags			login
//	@Summary		获取验证码
//	@Description	用户登录页获取验证码操作
//	@Produce		json
//	@Param			number	query		string	true	"手机号"
//	@Response		200,500	{object}	tools.HttpCode
//	@Router			/sendcode [GET]
func Sendcode(c *gin.Context) {
	number, _ := c.GetQuery("number")
	regex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	isValid := regex.MatchString(number)
	if !isValid {
		fmt.Println("错误手机号")
		c.JSON(http.StatusBadRequest, HttpCode{
			Code:    NotFound,
			Message: "手机号错误！",
			Data:    nil,
		})
		return
	}
	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html
	client, _err := CreateClient(tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")), tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")))
	if _err != nil {
		fmt.Println(_err)
		return
	}
	// 随机生成4位验证码
	code1 := fmt.Sprintf("%04d", rand.Intn(10000))
	//fmt.Println("验证码:", code1)
	var redisClient *redis.Client = model.RedisConn
	err := redisClient.Set(c, "captcha", code1, 5*time.Minute).Err()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, HttpCode{
			Code:    OK,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String(number),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code1)),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}
		//fmt.Println(a)
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			fmt.Println(_err)
			return
		}
	}
	c.JSON(http.StatusOK, HttpCode{
		Code:    OK,
		Message: "发送验证码",
	})
	return
}

// Codelogin
//
//	@Summary		验证码登录
//	@Description	执行登录操作
//	@Tags			login
//	@Accept			multipart/form-data
//	@Param			yzm		formData	string	true	"验证码"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/codelogin [POST]
func Codelogin(c *gin.Context) {
	formCode := c.PostForm("yzm")
	var redisClient *redis.Client = model.RedisConn
	redisCode, _ := redisClient.Get(c, "captcha").Result()
	if redisCode != formCode {
		c.JSON(http.StatusOK, HttpCode{
			Code:    NotFound,
			Message: "验证码错误",
		})
		return
	}
	c.JSON(http.StatusOK, HttpCode{
		Code:    OK,
		Message: "验证码登陆成功",
		Data:    nil,
	})
	return
}
