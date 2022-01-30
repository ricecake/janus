package model

type Action struct {
	Context string
	Name    string
	//TODO add a description
}

func (this Action) TableName() string {
	return "action"
}
