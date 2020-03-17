package model

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
