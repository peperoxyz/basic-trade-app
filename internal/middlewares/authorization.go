package middlewares

import (
	"basic-trade-app/internal/database"
	"basic-trade-app/internal/models"
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