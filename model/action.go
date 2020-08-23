package model

type Action struct {
	Context string
	Name    string
}

func (this Action) TableName() string {
	return "action"
}
