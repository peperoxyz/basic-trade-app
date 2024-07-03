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

func UpdateProduct(ctx *gin.Context) {
	db := database.GetDB()

	userData := ctx.MustGet("userData").(jwt5.MapClaims) // get userData from decoded jwt
	contentType := helpers.GetContentType(ctx) // get content type from request header
	Product := models.Product{}

	productUUID := ctx.Param("productUUID") // get product uuid from parameter
	userID := uint(userData["id"].(float64)) // get the userID from decoded jwt

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Product)
	} else {
		ctx.ShouldBind(&Product)
	}

	// retreive existing product from database that the uuid same with the param
	var getProduct models.Product
	if err := db.Model(&getProduct).Where("uuid = ?", productUUID).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	// update the product struct with the retreived data
	Product.ID = uint(getProduct.ID)
	Product.AdminID = userID

	updateData := models.Product {
		Name: Product.Name,
		ImageUrl: Product.ImageUrl,
	}

	if err := db.Model(&Product).Where("uuid = ?", productUUID).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	// retreive the updated data
	var updatedProduct models.Product
	if err := db.Where("uuid = ?", productUUID).First(&updatedProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	// response with json updated data
	ctx.JSON(http.StatusOK, gin.H {
		"data": updatedProduct,
	})






}