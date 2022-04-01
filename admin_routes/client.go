package admin_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func listClients(c *gin.Context) {
	ctx, err := model.ListClients()
	if err != nil {
		c.AbortWithStatusJSON(500, "Internal error")
		return
	}
	c.JSON(200, ctx)
}

type InputClient struct {
	Context     string
	DisplayName string
	Description string
	ClientId    string
	BaseUri     string
	Secret      *string
}

func updateClient(c *gin.Context) {
	var ctx InputClient
	if err := c.ShouldBind(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	client, err := model.FindClientById(ctx.ClientId)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	log.Printf("%+v", ctx)

	client.Context = ctx.Context
	client.DisplayName = ctx.DisplayName
	client.Description = ctx.Description
	client.BaseUri = ctx.BaseUri
	if ctx.Secret != nil {
		client.SetSecret(*ctx.Secret)
	}

	if err := client.SaveChanges(); err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, ctx)
}

func createClient(c *gin.Context) {
	var ctx model.Context
	if err := c.ShouldBind(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := model.CreateContext(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, ctx)

}
