package model

type Role struct {
	Context int
}

type UserRole struct {
	User  int
	Group int
	Role  int
}

type RoleAction struct {
	Role   int
	Action int
}
