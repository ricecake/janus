package model

import "github.com/ricecake/karma_chameleon/util"

type Action struct {
	Context string
	Name    string
	//TODO add a description
}

func (this Action) TableName() string {
	return "action"
}

func ListActions() ([]Action, error) {
	db := util.GetDb()
	var ctx []Action
	err := db.Find(&ctx).Error
	return ctx, err
}
