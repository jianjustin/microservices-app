package main

import (
	"go-micro.dev/v4/logger"
	"log"
	"os"

	httpServer "github.com/go-micro/plugins/v4/server/http"
	"go-micro.dev/v4"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

var Logger = log.New(os.Stderr, "", log.LstdFlags)

func main() {
	srv := httpServer.NewServer(
		server.Name("hello"),
		server.Address(":30000"),
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	demo := &demo{}
	demo.InitRouter(router)

	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		log.Fatalln(err)
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(registry.NewRegistry()),
	)
	service.Init()
	if err := service.Run(); err != nil {
		logger.Error(err)
		return
	}
}

type demo struct{}

func (a *demo) InitRouter(router *gin.Engine) {
	router.POST("/demo", a.demo)
}

func (a *demo) demo(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "call go-micro v4 http server success"})
}
