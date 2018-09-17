package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
}

func (user *User) First(){
	if err := Mdb.First(user).Error; err != nil {
		panic(err)
	}
}

func (user *User) Insert(){
	Mdb.Save(user)
}

func (User) TableName() string{
	return "users"
}