package controllers

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/helpers"
	"basic-trade-app/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateVariant(ctx *gin.Context) {
	db := database.GetDB()
	contentType := ctx.ContentType()

	Variant := models.Variant{}

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Variant)
	} else {
		ctx.ShouldBind(&Variant)
	}

	// get productUUID from form data request body
	productUUID := ctx.PostForm("product_uuid") 
	if productUUID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product UUID is required"})
		return
	}

	// get the product ID from the UUID
	Product := models.Product{}
	if err := db.Where("uuid = ?", productUUID).First(&Product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product with that UUID not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		}
		return
	}

	// set the product ID and generate a new UUID for the Variant
	Variant.ProductID = Product.ID
	Variant.UUID = uuid.New().String()

	// save the Variant to the database
	err := db.Debug().Create(&Variant).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H {
		"data": Variant,
	})
}

// get all variants with search by name
func GetVariants(ctx *gin.Context) {
	db := database.GetDB()
	var Variants []models.Variant

    // get query param
	variantName := ctx.Query("variant_name")

	limitStr := ctx.DefaultQuery("limit", "0")   // limit default 10
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

    // prepare query to find variants
	query := db.Debug().Model(&models.Variant{})

    // apply name filter if name parameter is exist
	if variantName != "" {
		query = query.Where("variant_name LIKE ?", "%"+variantName+"%")
	}

    // get the total count of products for pagination
    var totalCount int64
    query.Count(&totalCount)

    // apply pagination only if limit is greater than 0
    if limit > 0 {
        query = query.Limit(limit).Offset(offset * limit)
    }

	err = query.Find(&Variants).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

   // calculate total pages
    totalPages := 0
    if limit > 0 {
        totalPages = int(totalCount) / limit
        if int(totalCount)%limit != 0 {
            totalPages++
        }
    }

    // calculate current page
    currentPage := 1
    if limit > 0 {
        currentPage = offset + 1
    }

    ctx.JSON(http.StatusOK, gin.H{
        "count":   len(Variants),
        "data":    Variants,
        "total":   totalCount,
        "meta": gin.H{
            "current_page": currentPage,
            "total_pages":  totalPages,
            "limit":        limit,
            "offset":       offset,
        },
        "success": true,
		"message": "Variants retreived successfully",
    })
}



func GetVariant(ctx *gin.Context) {
	db := database.GetDB()
	variantUUID := ctx.Param("variantUUID")

	Variant := models.Variant{}
	err := db.Where("UUID = ?", variantUUID).First(&Variant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H {
				"error": "Variant with that UUID not found",
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
		"message": "Variant retreived successfully",
		"data": Variant,
	})
}

func UpdateVariant(ctx *gin.Context) {
	db := database.GetDB()

	contentType := helpers.GetContentType(ctx) // get content type from request header
	Variant := models.Variant{}

	variantUUID := ctx.Param("variantUUID") // get variant uuid from parameter

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Variant)
	} else {
		ctx.ShouldBind(&Variant)
	}

	// retreive existing variant from database that the uuid same with the param
	var getVariant models.Variant
	if err := db.Model(&getVariant).Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	// update the variant struct with the retreived data
	updateData := models.Variant {
	VariantName: Variant.VariantName,
	Quantity: Variant.Quantity,
	}

	if err := db.Model(&Variant).Where("uuid = ?", variantUUID).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	// // retreive the updated data
	var updatedVariant models.Variant
	if err := db.Where("uuid = ?", variantUUID).First(&updatedVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H {
		"data": updatedVariant,
	})
}

func DeleteVariant(ctx *gin.Context) {
	db := database.GetDB()
	Variant := models.Variant{}
	variantUUID := ctx.Param("variantUUID")


	// convert uuid from string to uint
	err := db.Where("UUID = ?", variantUUID).First(&Variant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H {
				"error": "Variant with that UUID not found",
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
	err = db.Delete(&Variant).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": "Internal server error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Variant deleted successfully",
	})
}