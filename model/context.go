package model

import (
	log "github.com/sirupsen/logrus"

	"github.com/ricecake/karma_chameleon/util"
)

type Context struct {
	Code        string `gorm:"column:code;not null; primary_key"`
	Name        string `gorm:"column:name;not null"`
	Description string `gorm:"column:description; not null"`
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

func ListContexts() ([]Context, error) {
	db := util.GetDb()
	var ctx []Context
	err := db.Find(&ctx).Error
	return ctx, err
}

func (this *Context) SaveChanges() error {
	db := util.GetDb()

	err := db.Save(this).Error
	if err != nil {
		log.Errorf("Error while updating context: %s", err)
	}

	return err
}
