package middlewares

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/models"
	"basic-trade-app/internal/models/requests"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

// to set only user who create the product that can do changes about that product
func ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		productUUID := ctx.Param("productUUID")

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		userID := uint(userData["id"].(float64))

		var getProduct models.Product
		err := db.Select("admin_id").Where("uuid = ?", productUUID).First(&getProduct).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H {
				"error": err.Error(),
				"message": "Data not found",
			})
			return
		}

		if getProduct.AdminID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H {
				"error": "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}
		ctx.Next()
	}
}

func VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		userID := uint(userData["id"].(float64))

		// Check if param variantUUID exists or not
		variantUUID := ctx.Param("variantUUID")
		if variantUUID != "" {
			// Proceed with variant authorization using the variantUUID parameter

			// Retrieve the variant and its associated product
			var getVariant models.Variant
			if err := db.Select("product_id").Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   err.Error(),
					"message": "Variant not found",
				})
				return
			}

			// Retrieve the product to check the admin ID
			var getProduct models.Product
			if err := db.Select("admin_id").Where("id = ?", getVariant.ProductID).First(&getProduct).Error; err != nil {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   err.Error(),
					"message": "Product not found",
				})
				return
			}

			// Check if the user is authorized
			if getProduct.AdminID != userID {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "You are not allowed to access this data",
				})
				return
			}
		} else {
			// Proceed with variant authorization using the request body (assuming it contains ProductUUID)

			// Bind the request body to extract ProductUUID
			var req requests.VariantRequest
			if err := ctx.ShouldBind(&req); err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":   "Bad request",
					"message": err.Error(),
				})
				return
			}

			// Retrieve the product associated with the ProductUUID
			var getProduct models.Product
			if err := db.Select("admin_id").Where("uuid = ?", req.ProductUUID).First(&getProduct).Error; err != nil {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   "Not found",
					"message": "Product not found",
				})
				return
			}

			// Check if the user is authorized
			if getProduct.AdminID != userID {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "You are not allowed to access this data",
				})
				return
			}
		}

		// If all checks pass, proceed to the next middleware or handler
		ctx.Next()
	}
}


/**
func VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()

		// Bind the request body to extract productUUID
		var req requests.VariantRequest
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		userID := uint(userData["id"].(float64))

		productUUID := req.ProductUUID

		// Retrieve the product associated with the productUUID
		var getProduct models.Product
		err := db.Select("admin_id").Where("uuid = ?", productUUID).First(&getProduct).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Not found",
				"message": "Product not found",
			})
			return
		}

		if getProduct.AdminID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}
		ctx.Next()
	}
}

func VariantAuthorizationByParam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		variantUUID := ctx.Param("variantUUID")

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		userID := uint(userData["id"].(float64))

		// Get the product associated with the variant
		var getVariant models.Variant
		err := db.Select("product_id").Where("uuid = ?", variantUUID).First(&getVariant).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Variant not found",
			})
			return
		}

		// Get the product to check the admin ID
		var getProduct models.Product
		err = db.Select("admin_id").Where("id = ?", getVariant.ProductID).First(&getProduct).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Product not found",
			})
			return
		}

		if getProduct.AdminID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}
		ctx.Next()
	}
}
*/