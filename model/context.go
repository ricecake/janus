package model

import (
	log "github.com/sirupsen/logrus"

	"github.com/ricecake/janus/util"
)

type Context struct {
	Code string `gorm:"column:code;not null"`
	Name string `gorm:"column:name;not null"`
}

func (this Context) TableName() string {
	return "context"
}

func CreateContext(ctx *Context) error {
	db := util.GetDb()

	if ctx.Code == "" {
		ctx.Code = util.CompactUUID()
	}
	log.Info("Creating context ", ctx.Code)

	return db.Create(ctx).Error
}
