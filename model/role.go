package model

import "github.com/ricecake/karma_chameleon/util"

type Role struct {
	Context   string
	Name      string
	Automatic bool
}

type UserCliqueRole struct {
	Context string
	User    string
	Clique  string
	Role    string
}

type UserRole struct {
	Context string
	User    string
	Role    string
}

type RoleAction struct {
	Role   string
	Action string
}

func ListRoles() ([]Role, error) {
	db := util.GetDb()
	var ctx []Role
	err := db.Find(&ctx).Error
	return ctx, err
}
