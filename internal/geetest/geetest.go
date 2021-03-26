package geetest

import (
	"github.com/gin-gonic/gin"
	"shortlink/internal/config"
	"shortlink/internal/geetest/sdk"
)

type GeetestApi struct {
	digestmod   string
	GEETEST_ID  string
	GEETEST_KEY string
}

func NewGeetestApi(conf *config.Setting) *GeetestApi {

	return &GeetestApi{
		digestmod:   "md5",
		GEETEST_ID:  conf.Geetest.GEETEST_ID,
		GEETEST_KEY: conf.Geetest.GEETEST_KEY,
	}
}

func (self *GeetestApi) FirstRegister(c *gin.Context) {
	/*
		   必传参数
		       digestmod 此版本sdk可支持md5
		   自定义参数,可选择添加
			   user_id 客户端用户的唯一标识，确定用户的唯一性；作用于提供进阶数据分析服务，可在register和validate接口传入，不传入也不影响验证服务的使用；若担心用户信息风险，可作预处理(如哈希处理)再提供到极验
			   client_type 客户端类型，web：电脑上的浏览器；h5：手机上的浏览器，包括移动应用内完全内置的web_view；native：通过原生sdk植入app应用的方式；unknown：未知
			   ip_address 客户端请求sdk服务器的ip地址
	*/
	var result *sdk.GeetestLibResult
	gtLib := sdk.NewGeetestLib(self.GEETEST_ID, self.GEETEST_KEY)
	digestmod := self.digestmod
	userID := "test"
	params := map[string]string{
		"digestmod":   digestmod,
		"user_id":     userID,
		"client_type": "web",
		"ip_address":  "127.0.0.1",
	}
	result = gtLib.Register(digestmod, params)
	c.Header("Content-Type", "application/json;charset=UTF-8")
	c.String(200, result.Data)
}

func (self *GeetestApi) SecondValidate(c *gin.Context) {
	gtLib := sdk.NewGeetestLib(self.GEETEST_ID, self.GEETEST_KEY)
	challenge := c.PostForm(sdk.GEETEST_CHALLENGE)
	validate := c.PostForm(sdk.GEETEST_VALIDATE)
	seccode := c.PostForm(sdk.GEETEST_SECCODE)
	var result *sdk.GeetestLibResult
	result = gtLib.SuccessValidate(challenge, validate, seccode)
	if result.Status == 1 {
		c.JSON(200, gin.H{"result": "success", "version": sdk.VERSION})
	} else {
		c.JSON(200, gin.H{"result": "fail", "version": sdk.VERSION, "msg": result.Msg})
	}
}
