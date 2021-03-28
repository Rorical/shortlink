package server

import (
	"shortlink/internal/config"
	"shortlink/internal/shortlink"

	"shortlink/internal/hcaptcha"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	api     *shortlink.ShortLinkApi
	captcha *hcaptcha.HCaptcha
}

func NewServer() *Server {
	router := gin.Default()
	conf := config.Read()
	api := shortlink.NewApi(conf)
	hcaptcha := hcaptcha.NewhCaptcha(conf)
	return &Server{
		router:  router,
		api:     api,
		captcha: hcaptcha,
	}
}

func makeJsonResponse(code int, message string, data map[string]interface{}) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}

func (self *Server) Init() {
	self.router.POST("/api/shorten/:id", func(c *gin.Context) {
		capres := self.captcha.GinVerify(c)
		if capres {
			id := c.Param("id")
			url := c.DefaultPostForm("url", "")
			err := self.api.SetLink(id, url)
			if err != nil {
				if err == shortlink.ErrAlreadyExists {
					c.JSON(200, makeJsonResponse(
						1,
						"Error Id already used",
						gin.H{},
					))
				} else if err == shortlink.ErrIllegalParameters {
					c.JSON(200, makeJsonResponse(
						-1,
						"Error Illegal Parameters",
						gin.H{},
					))
				} else {
					c.JSON(500, makeJsonResponse(
						500,
						"Unknown Error Please Contact With Administrator",
						gin.H{},
					))
				}
				return
			}
			c.JSON(200, makeJsonResponse(
				0,
				"OK",
				gin.H{},
			))
		} else {
			c.JSON(200, makeJsonResponse(
				-1,
				"Captcha Verify Failed",
				gin.H{},
			))
		}

	})

	self.router.GET("/api/short/:id/exist", func(c *gin.Context) {
		id := c.Param("id")
		exist, err := self.api.IsLinkExist(id)
		if err != nil {
			if err == shortlink.ErrIllegalParameters {
				c.JSON(200, makeJsonResponse(
					-1,
					"Error Illegal Parameters",
					gin.H{},
				))
			} else {
				c.JSON(500, makeJsonResponse(
					500,
					"Unknown Error Please Contact With Administrator",
					gin.H{},
				))
			}
			return
		}
		c.JSON(200, makeJsonResponse(
			0,
			"OK",
			gin.H{"exist": exist},
		))
	})

	self.router.GET("/api/short/:id", func(c *gin.Context) {
		id := c.Param("id")
		url, err := self.api.GetLink(id)
		if err != nil {
			if err == shortlink.ErrIllegalParameters {
				c.JSON(200, makeJsonResponse(
					-1,
					"Error Illegal Parameters",
					gin.H{},
				))
			} else if err == shortlink.ErrDoesNotExists {
				c.JSON(404, makeJsonResponse(
					2,
					"Error Record Not Found",
					gin.H{},
				))
			} else {
				c.JSON(500, makeJsonResponse(
					500,
					"Unknown Error Please Contact With Administrator",
					gin.H{},
				))
			}
			return
		}
		c.JSON(200, makeJsonResponse(
			0,
			"OK",
			gin.H{"url": url},
		))
	})

	self.InitPages()
}

func (self *Server) InitPages() {
	self.router.LoadHTMLGlob("internal/pages/*")
	self.router.GET("/redirect/:id", func(c *gin.Context) {
		id := c.Param("id")
		url, err := self.api.GetLink(id)
		if err != nil {
			if err == shortlink.ErrIllegalParameters {
				c.HTML(400, "error.tmpl", gin.H{
					"error":  -1,
					"detail": "Error Illegal Parameters",
				})
			} else if err == shortlink.ErrDoesNotExists {
				c.HTML(404, "error.tmpl", gin.H{
					"error":  2,
					"detail": "Error Record Not Found",
				})
			} else {
				c.HTML(500, "error.tmpl", gin.H{
					"error":  500,
					"detail": "Unknown Error Please Contact With Administrator",
				})
			}
			return
		}
		c.HTML(200, "redirect.tmpl", gin.H{
			"URL": url,
		})
	})
}

func (self *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	self.router.Run(":8081")
}
