package models

type User struct {
	Id       uint32 `sql:"AUTO_INCREMENT"`
	UserName string `json:"user"`
	Password string `json:"password"`
}
