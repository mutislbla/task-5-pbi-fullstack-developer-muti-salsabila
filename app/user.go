package app

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" valid:"alphanum,required"`
	Email    string `json:"email" gorm:"unique;not null" valid:"email,required"`
	Password string `json:"password,omitempty" valid:"minstringlength(6),required"`
}

type LoginData struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type ExistingUser struct {
	Id       int
	Username string
	Email    string
	Password string
}

type UserData struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
