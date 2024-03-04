package models

import (
	"task5-pbi/app"
)

type Photo struct {
	app.Photo
	UserID int
}
