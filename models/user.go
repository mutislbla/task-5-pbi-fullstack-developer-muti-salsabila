package models

import (
	"task5-pbi/app"
)

type User struct {
	app.User
	Photos []Photo `gorm:"foreignKey:UserID;contstraint:OnDelete:CASCADE" json:"photos"`
}
