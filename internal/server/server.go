package server

import (
	"shortlink/internal/config"
	"shortlink/internal/shortlink"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	api    *shortlink.ShortLinkApi
}

func NewServer() *Server {
	router := gin.Default()
	conf := config.Read()
	api := shortlink.NewApi(conf)
	return &Server{
		router: router,
		api:    api,
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
	self.router.GET("/api/shorten/:id", func(c *gin.Context) {
		id := c.Param("id")
		url := c.DefaultQuery("url", "")
		err := self.api.SetLink(id, url)
		if err != nil {
			if err == shortlink.ErrAlreadyExists {
				c.JSON(200, makeJsonResponse(
					1,
					"Error Id already used",
					gin.H{},
				))
			} else if err == shortlink.ErrIllegalParameters {
				c.JSON(400, makeJsonResponse(
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
	})

	self.router.GET("/api/short/:id/exist", func(c *gin.Context) {
		id := c.Param("id")
		exist, err := self.api.IsLinkExist(id)
		if err != nil {
			if err == shortlink.ErrIllegalParameters {
				c.JSON(400, makeJsonResponse(
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
				c.JSON(400, makeJsonResponse(
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
}

func (self *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	self.router.Run(":8080")
}
