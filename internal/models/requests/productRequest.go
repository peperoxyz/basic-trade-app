package requests

import "mime/multipart"

// dibuat karena terdapat upload file

type ProductRequest struct {
	Name 	string 					`form:"name" binding:"required"`
	Image 	*multipart.FileHeader 	`form:"file"`
}