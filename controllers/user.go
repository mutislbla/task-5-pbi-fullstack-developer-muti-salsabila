package controllers

import (
	"net/http"
	"strconv"
	"task5-pbi/app"
	"task5-pbi/database"
	"task5-pbi/helpers"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResponseData struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Register(c *gin.Context) {
	var user app.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if _, err := govalidator.ValidateStruct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := database.DB.Where("email=?", user.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is already in use"})
		return
	}
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user.Password = string(hashedPassword)
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ResponseData{
		Message: "Registration successful",
		Data: app.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

func Login(c *gin.Context) {
	var body app.LoginData
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var user app.ExistingUser
	if err := database.DB.Table("users").Where("email=?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email or password incorrect"})
		return
	}
	if err := helpers.CheckHashedPassword(body.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email or password incorrect"})
		return
	}
	accessToken := helpers.GenerateToken(user.Id, c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"token":   accessToken,
	})
}

func UpdateUserById(c *gin.Context) {
	userReq, _ := c.Get("user")
	userIdStr, _ := c.Params.Get("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
	}
	var body app.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	body.ID = userReq.(app.User).ID
	var user app.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if user.ID != body.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "No Access", "token id": user.ID, "body id": body.ID})
		return
	}
	if body.Password != "" {
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
		body.Password = string(passwordHash)
	}
	if err := database.DB.Model(&user).Where("id = ?", userId).Updates(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResponseData{
		Message: "User updated successfully",
		Data:    user,
	})
}

func DeleteUserById(c *gin.Context) {
	userReq, _ := c.Get("user")
	userIdStr, _ := c.Params.Get("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
	}
	var user app.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if user.ID != userReq.(app.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Not Permitted"})
		return
	}
	if err := database.DB.Delete(&user, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "User deleted successfully",
	})
}
