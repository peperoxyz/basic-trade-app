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