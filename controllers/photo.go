package controllers

import (
	"net/http"
	"strconv"
	"task5-pbi/app"
	"task5-pbi/database"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func AddPhoto(c *gin.Context) {
	user, _ := c.Get("user")
	var body app.Photo
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}
	body.UserID = user.(app.User).ID
	if err := database.DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ResponseData{
		Message: "Photo added successfully",
		Data:    body,
	})
}

func GetAllPhotos(c *gin.Context) {
	var photos []app.Photo
	if err := database.DB.Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResponseData{
		Message: "Get all photos successfully",
		Data:    photos,
	})
}

func UpdatePhotoById(c *gin.Context) {
	user, _ := c.Get("user")
	var photo app.Photo
	photoIdStr, _ := c.Params.Get("photoId")
	photoId, err := strconv.Atoi(photoIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid photo id"})
	}
	var body app.Photo
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}
	body.UserID = user.(app.User).ID
	if err := database.DB.Where("id = ?", photoId).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": err.Error()})
		return
	}
	if photo.UserID != body.UserID {
		c.JSON(http.StatusForbidden, gin.H{"Message": "Not Permitted"})
		return
	}
	if err := database.DB.Model(&photo).Where("id = ?", photoId).Updates(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResponseData{
		Message: "Photo updated successfully",
		Data:    photo,
	})
}

func DeletePhotoById(c *gin.Context) {
	user, _ := c.Get("user")
	photoIdStr, _ := c.Params.Get("photoId")
	photoId, err := strconv.Atoi(photoIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid User ID"})
	}
	var photo app.Photo
	if err := database.DB.Where("id = ?", photoId).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": err.Error()})
		return
	}
	if photo.UserID != user.(app.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"Message": "Not Permitted"})
		return
	}
	if err := database.DB.Delete(&photo, photoId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Photo deleted successfully",
	})
}
