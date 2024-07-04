package controllers

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/helpers"
	"basic-trade-app/internal/models"
	"basic-trade-app/internal/models/requests"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// create Product with upload image
func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()

	var productReq requests.ProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	// extract filename without extention
	fileName := helpers.RemoveExtension(productReq.Image.Filename)

	uploadResult, err := helpers.UploadFile(productReq.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}
	
	Product := models.Product{
		Name: productReq.Name,
		ImageUrl: uploadResult,
	}

	// from authentication return
	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	AdminID := uint(userData["id"].(float64))


	Product.AdminID = AdminID
	newUUID := uuid.New()
	Product.UUID = newUUID.String()

	// proses upload image
	err = db.Debug().Create(&Product).Error
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

// function get products with search by name and pagination
func GetProducts(ctx *gin.Context) {
	db := database.GetDB()
	var Products []models.Product

	// get query param
	name := ctx.Query("name")

	limitStr := ctx.DefaultQuery("limit", "10")   // limit default 10
	offsetStr := ctx.DefaultQuery("offset", "0")  // offset default 0

	// konversi limit dan offset menjadi integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad request",
			"message": "parameter limit tidak valid",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad request",
			"message": "parameter offset tidak valid",
		})
		return
	}

	// prepare query to find products
	query := db.Debug().Model(&models.Product{})

	// apply name filter if name parameter is exist
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// preload variants data
	query = query.Preload("Variants")

	// terapkan pagination
	query = query.Limit(limit).Offset(offset)

	// execute query to find the products wanted
	err = query.Find(&Products).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": "Internal server error",
			"message": err.Error(),
		})
		return
	}
	
	Count := len(Products)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    Products,
		"count": Count,
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

	// bind req data to productReq struct
	var productReq requests.ProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	// extract filename without extention
	fileName := helpers.RemoveExtension(productReq.Image.Filename)
	
	// upload the file and get the url
	uploadResult, err := helpers.UploadFile(productReq.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	// prepare the updated product data
	Product := models.Product{
		Name: productReq.Name,
		ImageUrl: uploadResult,
	}

	// get the user data from JWT token
	userData := ctx.MustGet("userData").(jwt5.MapClaims) // get userData from decoded jwt
	userID := uint(userData["id"].(float64)) // get the userID from decoded jwt

	// get product UUID from parameter
	productUUID := ctx.Param("productUUID") // get product uuid from parameter

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

func GetProduct(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product retreived successfully",
		"data": Product,
	})
}
