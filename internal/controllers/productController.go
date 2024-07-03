package controllers

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/helpers"
	"basic-trade-app/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()
	
	// from authentication return
	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	contentType := helpers.GetContentType(ctx)

	Product := models.Product{}
	AdminID := uint(userData["id"].(float64))

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Product)
	} else {
		ctx.ShouldBind(&Product)
	}

	Product.AdminID = AdminID
	newUUID := uuid.New()
	Product.UUID = newUUID.String()

	err := db.Debug().Create(&Product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H {
		"data": Product,
	})
}

func GetProducts(ctx *gin.Context) {
	db := database.GetDB()
	// Products := models.Product{}
	var Products []models.Product
	
	// Find all products
	// err := db.Debug().Preload("Admin").Find(&Products).Error
	err := db.Debug().Find(&Products).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    Products,
	})
}

func DeleteProduct(ctx *gin.Context) {
	db := database.GetDB()
	productUUID := ctx.Param("productUUID")

	// convert uuid from string to uint
	Product := models.Product{}
	err := db.Where("UUID = ?", productUUID).First(&Product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H {
				"error": "Product with that UUID not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"error": "Internal server error",
				"message": err.Error(),
			})
		}
		return
	}

	// if no error, delete the product
	err = db.Delete(&Product).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": "Internal server error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"data": Product,
	})
}