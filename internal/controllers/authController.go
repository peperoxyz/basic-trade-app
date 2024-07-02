package controllers

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/helpers"
	"basic-trade-app/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRegister(ctx *gin.Context) {
	db := database.GetDB() // connect to db
	contentType := helpers.GetContentType(ctx) // get content type
	Admin := models.Admin{}

	// define if the request content type formData, use ctx.ShouldBind. but if it's application/json, use ctx.ShouldBindJSON
	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Admin)
	} else {
		ctx.ShouldBind(&Admin)
	}

	/*
	// generate new UUID
	newUUID := uuid.New()
	Admin.UUID = newUUID.String() // set the generated uuid as the user's uuid field
	*/

	err := db.Debug().Create(&Admin).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error" : "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": Admin,
	})
}

func AdminLogin(ctx *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(ctx)
	Admin := models.Admin{}
	var password string

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Admin)
	} else {
		ctx.ShouldBind(&Admin)
	}

	password = Admin.Password // inputtan user

	// find admin with inputted email
	err := db.Debug().Where("email = ?", Admin.Email).Take(&Admin).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H {
			"error": "Unauthorized",
			"message": "Invalid email",
		})
		return
	}


	// if email found, compare the hashed pass to the inputted pass
	comparePass := helpers.ComparePass([]byte(Admin.Password), []byte(password))
	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H {
			"error": "Unauthorized",
			"message": "Invalid password",
		})
		return
	}

	// if no problems found with the email and password, generate the token
	token := helpers.GenerateToken(Admin.ID, Admin.UUID, Admin.Email)

	ctx.JSON(http.StatusOK, gin.H {
		"token": token,
		"message": "Successfully logged in",
	})


}