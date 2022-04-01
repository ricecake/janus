package model

import "github.com/ricecake/karma_chameleon/util"

type Clique struct {
	Context string
	Name    string
}

func ListGroups() ([]Clique, error) {
	db := util.GetDb()
	var ctx []Clique
	err := db.Find(&ctx).Error
	return ctx, err
}
